package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"hash/crc32"
	"log"
	"sort"

	"oddstream.games/gosold/cardid"
	"oddstream.games/gosold/util"
)

// Baize holds the state of the baize, piles and cards therein.
// Baize is exported from this package because it's used to pass between light and dark.
// LIGHT should see a Baize object as immutable, hence the unexported fields and getters.
type Baize struct {
	dark         *dark  // link back to dark for statistics
	variant      string // needed by stats (could fix this)
	pilesToCheck []*Pile

	// members needed by solver
	script    scripter
	cardCount int
	piles     []*Pile // needed by LIGHT to display piles and cards
	recycles  int     // needed by LIGHT to determine Stock rune
	percent   int     // needed by LIGHT to display in status bar

	// members specific to solver
	depth      int
	parent     *Baize
	children   []*Baize
	crc        uint32
	tappedCard cardid.CardID

	// members that are needed by LIGHT
	bookmark  int // needed by LIGHT to grey out goto bookmark menu item
	moves     int // count of all available card moves
	fmoves    int // count of available moves to foundation
	undoStack []*savableBaize
	fnNotify  func(BaizeEvent, any)
	BaizeSettings
}

type BaizeEvent int

const (
	ChangedEvent BaizeEvent = iota
	LabelEvent
	WonEvent
	LostEvent
	MessageEvent
)

type BaizeSettings struct {
	PowerMoves, SafeCollect bool
}

// NewBaize creates a new Baize object
func (d *dark) NewBaize(variant string, fnNotify func(BaizeEvent, any)) (*Baize, error) {
	var script scripter
	var ok bool
	if script, ok = variants[variant]; !ok {
		return nil, errors.New("unknown variant " + variant)
	}
	b := &Baize{dark: d, variant: variant, script: script, fnNotify: fnNotify}
	b.script.SetBaize(b)
	b.script.BuildPiles()
	// BuildPiles() must not create or move any cards
	// so fill and shuffle Stock here
	{
		stock := b.script.Stock()
		b.cardCount = stock.fill(b.script.Packs(), b.script.Suits())
		stock.shuffle()
	}
	b.pilesToCheck = []*Pile{}
	b.pilesToCheck = append(b.pilesToCheck, b.script.Foundations()...)
	b.pilesToCheck = append(b.pilesToCheck, b.script.Tableaux()...)
	b.pilesToCheck = append(b.pilesToCheck, b.script.Cells()...)
	b.pilesToCheck = append(b.pilesToCheck, b.script.Discards()...)
	if b.script.Waste() != nil {
		// in Go 1.19, append will add a nil
		// in Go 1.17, nil was not appended?
		b.pilesToCheck = append(b.pilesToCheck, b.script.Waste())
	}

	b.script.StartGame()
	// NOT calling afterChange() here; we don't want the notification sent
	b.undoPush()
	b.findTapTargets()
	b.percent = b.percentComplete()
	return b, nil
}

// Baize public interface ////////////////////////////////////////////////////////////

func (b *Baize) Bookmark() int {
	return b.bookmark
}

func (b *Baize) Complete() bool {
	return b.script.Complete()
}

func (b *Baize) Conformant() bool {
	for _, p := range b.piles {
		if !p.vtable.Conformant() {
			return false
		}
	}
	return true
}

func (b *Baize) LoadPosition() (bool, error) {
	if b.bookmark == 0 || b.bookmark > len(b.undoStack) {
		return false, errors.New("No bookmark")
	}
	if b.Complete() {
		return false, errors.New("Cannot undo a completed game") // otherwise the stats can be cooked
	}
	var sav *savableBaize
	var ok bool
	for len(b.undoStack)+1 > b.bookmark {
		sav, ok = b.undoPop()
		if !ok {
			log.Panic("error popping from undo stack")
		}
	}
	b.updateFromSavable(sav)
	b.undoPush()
	b.afterChange()
	return true, nil
}

func (b *Baize) PercentComplete() int {
	return b.percent
}

// Piles returns the slice of Piles
func (b *Baize) Piles() []*Pile {
	return b.piles
}

// PileTapped called by client when a pile - usually stock - has been tapped
func (b *Baize) PileTapped(pile *Pile) bool {
	cardsMoved := false
	oldCRC := b.calcCRC()
	b.script.PileTapped(pile)
	if b.calcCRC() != oldCRC {
		b.afterUserMove()
		cardsMoved = true
		b.fnNotify(ChangedEvent, nil)
	}
	return cardsMoved
}

func (b *Baize) Recycles() int {
	return b.recycles
}

func (b *Baize) CardColors() int {
	return b.script.CardColors()
}

func (b *Baize) NewDeal() (bool, error) {
	// a virgin game has one state on the undo stack
	if len(b.undoStack) > 1 && !b.Complete() {
		percent := b.PercentComplete()
		b.dark.stats.recordLostGame(b.variant, percent)
		b.fnNotify(LostEvent, nil)
	}

	b.reset()

	// the cards are all over the place
	// we could either recall them all to the Stock
	// or just delete them and make fresh ones
	// favour the former because the cards lerp better

	// for _, p := range b.piles {
	// 	p.cards = p.cards[:0]
	// }
	// b.cardCount = b.script.Stock().fill(b.script.Packs(), b.script.Suits())

	stock := b.script.Stock()
	for _, p := range b.piles {
		if p == stock {
			continue
		}
		stock.cards = append(stock.cards, p.cards...)
		p.cards = p.cards[:0]
	}
	for _, c := range stock.cards {
		c.pile = stock
	}
	if len(stock.cards) != b.cardCount {
		log.Panic("the number of cards in the stock is incorrect")
	}

	b.script.Stock().shuffle()
	b.script.StartGame()
	b.undoPush()
	b.afterChange()
	return true, nil
}

func (b *Baize) RestartDeal() (bool, error) {
	if b.Complete() {
		return false, errors.New("Cannot restart a completed game") // otherwise the stats can be cooked
	}
	var sav *savableBaize
	var ok bool
	for len(b.undoStack) > 0 {
		sav, ok = b.undoPop()
		if !ok {
			log.Panic("error popping from undo stack")
		}
	}
	b.updateFromSavable(sav)
	b.bookmark = 0 // do this AFTER UpdateFromSavable
	b.undoPush()
	b.afterChange()
	return true, nil
}

func (b *Baize) Load() {
	// pearl from the mudbank:
	// don't do a crc check here; send the change notify in all cases
	if !NoLoad {
		b.load()
		b.fnNotify(ChangedEvent, nil)
	}
}

func (b *Baize) Save() {
	if !NoSave {
		b.save()
	}
}

// SavePosition sets the bookmark to the current baize position
func (b *Baize) SavePosition() (bool, error) {
	if b.Complete() {
		return false, errors.New("Cannot bookmark a completed game") // otherwise the stats can be cooked
	}
	b.bookmark = len(b.undoStack)
	sb := b.undoPeek()
	sb.Bookmark = b.bookmark
	sb.Recycles = b.recycles
	return true, nil
}

func (b *Baize) SetSettings(settings BaizeSettings) {
	b.BaizeSettings = settings
}

// TailDragged called by client when a tail of cards has been dragged from one pile to another.
// If this func returns false, the client should animate the tail back to where it came from
// and toast any error message.
func (b *Baize) TailDragged(src *Pile, tail []*Card, dst *Pile) (bool, error) {
	if src == dst {
		return false, nil // put the tail back, but don't make a fuss about it
	}
	var ok bool
	var err error
	if ok, err = src.canMoveTail(tail); !ok {
		return false, err
	} else {
		if ok, err = dst.vtable.CanAcceptTail(tail); !ok {
			return false, err
		} else {
			if ok, err = b.script.TailMoveError(tail); !ok {
				return false, err
			} else {
				oldCRC := b.calcCRC()
				if len(tail) == 1 {
					moveCard(src, dst)
				} else {
					moveTail(tail[0], dst)
				}
				if b.calcCRC() != oldCRC {
					b.afterUserMove()
					b.fnNotify(ChangedEvent, nil)
				}
			}
		}
	}
	return true, nil
}

func (b *Baize) CardDragged(src *Pile, card *Card, dst *Pile) (bool, error) {
	// tail := card.owningPile.makeTail(card)
	tail := src.makeTail(card)
	return b.TailDragged(src, tail, dst)
}

// TailTapped called by client when a card/tail has been tapped.
// returns true if cards have been moved.
func (b *Baize) TailTapped(tail []*Card, nTarget int) bool {
	cardsMoved := false
	oldCRC := b.calcCRC()
	b.script.TailTapped(tail, nTarget)
	if b.calcCRC() != oldCRC {
		b.afterUserMove()
		cardsMoved = true
		if b.fnNotify != nil {
			b.fnNotify(ChangedEvent, nil)
		}
	}
	return cardsMoved
}

func (b *Baize) CardTapped(card *Card, nTarget int) bool {
	tail := card.pile.makeTail(card)
	return b.TailTapped(tail, nTarget)
}

func (b *Baize) Undo() (bool, error) {
	if len(b.undoStack) < 2 {
		return false, errors.New("Nothing to undo")
	}
	if b.Complete() {
		return false, errors.New("Cannot undo a completed game") // otherwise the stats can be cooked
	}
	_, ok := b.undoPop() // removes current state
	if !ok {
		log.Panic("error popping current state from undo stack")
	}

	sav, ok := b.undoPop() // removes previous state for examination
	if !ok {
		log.Panic("error popping second state from undo stack")
	}
	b.updateFromSavable(sav)
	b.undoPush()
	b.afterChange()
	return true, nil
}

func (b *Baize) UndoStackSize() int {
	return len(b.undoStack)
}

func (b *Baize) Moves() (int, int) {
	return b.moves, b.fmoves
}

// StockLen returns number of cards in Stock, or -1 if Stock is hidden.
func (b *Baize) StockLen() int {
	if b.script.Stock().Hidden() {
		return -1
	}
	return b.script.Stock().Len()
}

// WasteLen returns number of cards in Waste, or -1 if there is no Waste.
func (b *Baize) WasteLen() int {
	if b.script.Waste() == nil {
		return -1
	}
	return b.script.Waste().Len()
}

// collectFromPile is a helper function for Collect()
func (b *Baize) collectFromPile(pile *Pile) int {
	if pile == nil {
		return 0
	}
	var cardsMoved int = 0
	for _, fp := range b.script.Foundations() {
		for {
			var card *Card = pile.peek()
			if card == nil {
				return cardsMoved
			}
			ok, _ := fp.vtable.CanAcceptTail([]*Card{card})
			if !ok {
				break // done with this foundation, try another
			}
			if b.SafeCollect {
				if ok, safeOrd := b.doingSafeCollect(); ok {
					if card.Ordinal() > safeOrd {
						// can't toast here, collect all will create a lot of toasts
						// TheGame.UI.Toast("Glass", fmt.Sprintf("Unsafe to collect %s", card.String()))
						break // done with this foundation, try another
					}
				}
			}
			moveCard(pile, fp)
			b.afterUserMove() // does an undoPush()
			cardsMoved += 1
		}
	}
	return cardsMoved
}

// Collect should be exactly the same as the user tapping repeatedly on the
// waste, cell, reserve and tableau piles.
// nb there is no collecting to discard piles, they are optional and presence of
// cards in them does not signify a complete game.
func (b *Baize) Collect() {
	var totalCardsMoved int = 0
	for {
		var cardsMoved int = b.collectFromPile(b.script.Waste())
		for _, pile := range b.script.Cells() {
			cardsMoved += b.collectFromPile(pile)
		}
		for _, pile := range b.script.Reserves() {
			cardsMoved += b.collectFromPile(pile)
		}
		for _, pile := range b.script.Tableaux() {
			cardsMoved += b.collectFromPile(pile)
		}
		if cardsMoved == 0 {
			break
		}
		totalCardsMoved += cardsMoved
	}
	if totalCardsMoved > 0 {
		b.fnNotify(ChangedEvent, nil)
	}
}

func (b *Baize) Wikipedia() string {
	return b.script.Wikipedia()
}

// Baize private functions

func (b *Baize) reset() {
	b.undoStack = nil
	b.bookmark = 0
}

// percentComplete used to display in status bar, and as postive progress in solver
func (b *Baize) percentComplete() int {
	var pairs, unsorted, percent int
	for _, p := range b.piles {
		if p.Len() > 1 {
			pairs += p.Len() - 1
		}
		unsorted += p.vtable.unsortedPairs()
	}
	percent = (int)(100.0 - util.MapValue(float64(unsorted), 0, float64(pairs), 0.0, 100.0))
	return percent
}

func (b *Baize) calcCRC() uint32 {
	/*
		var crc uint = 0xFFFFFFFF
		var mask uint
		for _, p := range b.piles {
			crc = crc ^ uint(p.Len())
			for j := 7; j >= 0; j-- {
				mask = -(crc & 1)
				crc = (crc >> 1) ^ (0xEDB88320 & mask)
			}
		}
		crc = ^crc // bitwise NOT
		return uint32(crc)
	*/
	var arr []byte // = []byte{byte(b.recycles)}
	for _, p := range b.piles {
		arr = append(arr, byte(p.Len()))
		for _, c := range p.cards {
			arr = append(arr, byte(c.id.Pack()))
			arr = append(arr, byte(c.id.Suit()))
			arr = append(arr, byte(c.id.Ordinal()))
		}
	}
	return crc32.ChecksumIEEE(arr)
}

func (b *Baize) addPile(pile *Pile) {
	b.piles = append(b.piles, pile)
}

func (b *Baize) setRecycles(recycles int) {
	if b.recycles != recycles {
		b.recycles = recycles
		b.fnNotify(LabelEvent, b.script.Stock())
	}
}

func (b *Baize) afterChange() {
	b.findTapTargets()
	b.percent = b.percentComplete()
	b.fnNotify(ChangedEvent, nil)
}

func (b *Baize) afterUserMove() {
	b.script.AfterMove()
	b.undoPush()
	b.afterChange()
	if b.Complete() {
		b.dark.stats.recordWonGame(b.variant, len(b.undoStack)-1)
		b.fnNotify(WonEvent, nil)
	}
}

func (b *Baize) replaceUndoStack(undoStack []*savableBaize) {
	b.undoStack = undoStack
	sav := b.undoPeek()
	b.updateFromSavable(sav)
	// NOT calling undoPush() here
	// because undo stack is replaced by the loaded one
	b.afterChange()
}

func (b *Baize) calcPowerMoves(pDraggingTo *Pile) int {
	// (1 + number of empty freecells) * 2 ^ (number of empty columns)
	// see http://ezinearticles.com/?Freecell-PowerMoves-Explained&id=104608
	// and http://www.solitairecentral.com/articles/FreecellPowerMovesExplained.html
	var emptyCells, emptyCols int
	for _, p := range b.piles {
		if p.Empty() {
			switch p.vtable.(type) {
			case *Cell:
				emptyCells++
			case *Tableau:
				if p.Label() == "" && p != pDraggingTo {
					// 'If you are moving into an empty column, then the column you are moving into does not count as empty column.'
					emptyCols++
				}
			}
		}
	}
	// 2^1 == 2, 2^0 == 1, 2^-1 == 0.5
	n := (1 + emptyCells) * util.Pow(2, emptyCols)
	// println(emptyCells, "emptyCells,", emptyCols, "emptyCols,", n, "powerMoves")
	return n
}

// DoingSafeCollect returns true (if we are doing safe collect)
// and the safe ordinal to collect next
func (b *Baize) doingSafeCollect() (bool, int) {
	if !b.script.SafeCollect() {
		return false, 0
	}
	var fs []*Pile = b.script.Foundations()
	if fs == nil {
		return false, 0
	}
	var f0 *Pile = fs[0]
	if f0 == nil {
		return false, 0
	}
	if f0.Label() != "A" {
		return false, 0 // eg Duchess
	}
	var lowest int = 99
	for _, f := range fs {
		if f.Empty() {
			// it's okay to collect aces and twos to start with
			return true, 2
		}
		var card *Card = f.peek()
		if card.Ordinal() < lowest {
			lowest = card.Ordinal()
		}
	}
	return true, lowest + 1
}

// foreachCard applys a function to each card
func (b *Baize) foreachCard(fn func(*Card)) {
	for _, p := range b.piles {
		for _, c := range p.cards {
			fn(c)
		}
	}
}

func (b *Baize) findAllMovableTails2() [][]*Card {
	var tails [][]*Card
	for _, p := range b.piles {
		if ptails := p.vtable.MovableTails2(); ptails != nil {
			tails = append(tails, ptails...)
		}
	}
	return tails
}

func (b *Baize) findTargetsForAllMovableTails2(tails [][]*Card) {

	for _, tail := range tails {
		// we already know this tail is movable, both at pile-type and script level
		headCard := tail[0]
		src := headCard.pile
		for _, dst := range b.pilesToCheck {
			if dst == src {
				continue
			}
			if ok, _ := dst.vtable.CanAcceptTail(tail); ok {
				// moving an full tail from one pile to another empty pile of the same type is pointless
				// eg Cell to Cell or Tableau to Tableau
				if dst.Len() == 0 && src.Len() == len(tail) && src.label == dst.label && src.category == dst.category {
					continue
				}
				// filter out case of, for example, moving a single card from
				// tableau to any of four different empty cells; just record one
				if dst.Len() == 0 && len(tail) == 1 {
					if dst.category == "Cell" || dst.category == "Tableau" || dst.category == "Foundation" {
						// if slices.ContainsFunc(headCard.tapTargets, func(tt tapTarget) bool { return tt.dst.category == dst.category }) {
						// 	continue
						// }
						contains := false
						for _, tt := range headCard.tapTargets {
							if tt.dst.category == dst.category {
								contains = true
								break
							}
						}
						if contains {
							continue
						}
					}
				}
				var weight int16
				switch dst.vtable.(type) {
				case *Cell:
					weight = 1
				case *Tableau:
					if dst.Empty() {
						if dst.Label() != "" {
							weight = 2
						} else {
							weight = 1
						}
					} else if dst.peek().Suit() == headCard.Suit() {
						// Simple Simon, Spider
						weight = 3
					} else {
						weight = 2
					}
				case *Foundation, *Discard:
					// moves to Foundation get priority when card is tapped
					weight = 4
				default:
					weight = 1
				}
				headCard.tapTargets = append(headCard.tapTargets, tapTarget{dst: dst, weight: weight})
			}
		}
	}
}

func (b *Baize) sortTapTargets() {
	for _, p := range b.piles {
		for _, c := range p.cards {
			if c.tapTargets == nil {
				continue
			}
			// sort so highest weight comes first
			sort.Slice(c.tapTargets, func(i, j int) bool { return c.tapTargets[i].weight > c.tapTargets[j].weight })
		}
	}
}

func (b *Baize) countMoves() {
	b.moves, b.fmoves = 0, 0

	if !b.script.Stock().Hidden() {
		if b.script.Stock().Empty() {
			if b.Recycles() > 0 {
				b.moves++
			}
		} else {
			// games like Agnes B (with a Spider-like stock) need to report an available move
			// so we can't do this:
			// card := b.script.Stock().peek()
			// card.destinations = b.FindHomesForTail([]*Card{card})
			// b.moves += len(card.destinations)
			b.moves += 1
		}
	}

	for _, p := range b.piles {
		for _, c := range p.cards {
			if c.tapTargets == nil {
				continue
			}
			for _, tt := range c.tapTargets {
				b.moves++
				if _, ok := tt.dst.vtable.(*Foundation); ok {
					b.fmoves++
				}
			}
		}
	}

}

func (b *Baize) findTapTargets() {

	// TODO save this in Baize and calculate after script.BuildPiles()
	b.foreachCard(func(c *Card) { c.tapTargets = nil })
	var tails [][]*Card = b.findAllMovableTails2()
	b.findTargetsForAllMovableTails2(tails) // adds tapTargets to movable cards
	b.sortTapTargets()
	b.countMoves()
}
