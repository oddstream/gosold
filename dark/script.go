package dark

import (
	"log"
)

// scripter defines the interface that variant-specific 'scripts' must supply,
// albeit several will be supplied by the embedded ScriptBase struct.
// TODO for the moment, methods are published.
type scripter interface {
	SetBaize(*Baize)
	Reset()
	BuildPiles()
	StartGame()
	AfterMove()

	TailMoveError([]*Card) (bool, error)
	TailAppendError(*Pile, []*Card) (bool, error)

	TailTapped([]*Card)
	PileTapped(*Pile)

	Cells() []*Pile
	Discards() []*Pile
	Foundations() []*Pile
	Reserves() []*Pile
	Stock() *Pile
	Tableaux() []*Pile
	Waste() *Pile
	Wastes() []*Pile

	Complete() bool
	Wikipedia() string
	CardColors() int
	Packs() int
	Suits() int

	Script() string
}

type scriptBase struct {
	baize        *Baize
	cells        []*Pile
	discards     []*Pile
	foundations  []*Pile
	reserves     []*Pile
	stock        *Pile
	tableaux     []*Pile
	wastes       []*Pile
	wikipedia    string
	cardColors   int
	packs, suits int
	script       string // empty for a builtin
	// could add suitFilter
}

// Fallback/default methods for a scripter interface //////////////////////////

func (sb *scriptBase) SetBaize(b *Baize) {
	sb.baize = b
}

// Reset is needed when changing variants that use the same class
func (sb *scriptBase) Reset() {
	sb.cells = nil
	sb.discards = nil
	sb.foundations = nil
	sb.reserves = nil
	sb.stock = nil
	sb.tableaux = nil
	sb.wastes = nil
}

// no default for BuildPiles

// no default for StartGame

func (sb scriptBase) AfterMove() {}

func (sb scriptBase) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmpFunc)
}

func (sb scriptBase) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail)
	}
	return dst.appendCmpFunc(dyad{dst.peek(), tail[0]})
}

func (sb scriptBase) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if len(tail) == 1 && pile == sb.Stock() && !sb.Stock().Hidden() && len(sb.Wastes()) == 1 {
		moveCard(sb.Stock(), sb.Waste())
	} else {
		pile.vtable.tailTapped(tail)
	}
}

func (sb scriptBase) PileTapped(*Pile) {}

func (sb scriptBase) Cells() []*Pile {
	return sb.cells
}

func (sb scriptBase) Discards() []*Pile {
	return sb.discards
}

func (sb scriptBase) Foundations() []*Pile {
	return sb.foundations
}

func (sb scriptBase) Reserves() []*Pile {
	return sb.reserves
}

func (sb scriptBase) Reserve() *Pile {
	if sb.reserves == nil {
		return nil
	}
	return sb.reserves[0]
}

func (sb scriptBase) Stock() *Pile {
	return sb.stock
}

func (sb scriptBase) Tableaux() []*Pile {
	return sb.tableaux
}

func (sb scriptBase) Waste() *Pile {
	if sb.wastes == nil {
		return nil
	}
	return sb.wastes[0]
}

func (sb scriptBase) Wastes() []*Pile {
	return sb.wastes
}

func (sb scriptBase) Complete() bool {
	if len(sb.discards) > 0 {
		for _, tab := range sb.tableaux {
			switch len(tab.cards) {
			case 0:
				// that's fine
			case 13:
				// TODO BUG this is currently wrong
				// as the logic allows a pile of mixed suit cards
				// to be complete (eg in Simple Simon, Spider Two Suits)
				// (lsol does not have this problem)
				if !tab.vtable.conformant() {
					return false
				}
			default:
				return false
			}
		}
		return true
	} else {
		// In Bisley, there may be <13 cards in a Foundation
		var n = 0
		for _, f := range sb.foundations {
			n += len(f.cards)
		}
		return n == sb.baize.numberOfCards()
	}
}

func (sb scriptBase) Wikipedia() string {
	if sb.wikipedia == "" { // uninitialized default
		return "https://en.wikipedia.org/wiki/Patience_(game)"
	} else {
		return sb.wikipedia
	}
}

func (sb scriptBase) CardColors() int {
	if sb.cardColors == 0 {
		return 2
	}
	return sb.cardColors
}

func (sb scriptBase) Packs() int {
	if sb.packs == 0 {
		return 1
	}
	return sb.packs
}

func (sb scriptBase) Suits() int {
	if sb.suits == 0 {
		return 4
	}
	return sb.suits
}

func (sb scriptBase) Script() string {
	return sb.script
}

// useful generic game library of functions ///////////////////////////////////

func anyCardsProne(cards []*Card) bool {
	for _, c := range cards {
		if c.Prone() {
			return true
		}
	}
	return false
}

// moveCard moves the top card from src to dst
func moveCard(src *Pile, dst *Pile) *Card {
	if c := src.pop(); c != nil {
		dst.push(c)
		src.flipUpExposedCard()
		return c
	}
	return nil
}

// moveTail moves all the cards from card downwards onto dst
func moveTail(card *Card, dst *Pile) {
	var src *Pile = card.owner()
	tmp := make([]*Card, 0, len(src.cards))
	// pop cards from src upto and including the head of the tail, onto a tmp stack
	for {
		var c *Card = src.pop()
		if c == nil {
			log.Panicf("MoveTail could not find %s", card)
		}
		tmp = append(tmp, c)
		if c == card {
			break
		}
	}
	// pop cards from the tmp stack and push onto dst
	if len(tmp) > 0 {
		for len(tmp) > 0 {
			var c *Card = tmp[len(tmp)-1]
			tmp = tmp[:len(tmp)-1]
			dst.push(c)
		}
		src.flipUpExposedCard()
	}
}

// populateWasteFromStock move n cards from stock to waste if waste is empty
func (sb *scriptBase) populateWasteFromStock(n int) {
	if sb.Waste() != nil {
		if sb.Waste().Len() == 0 {
			for i := 0; i < n; i++ {
				moveCard(sb.stock, sb.Waste())
			}
		}
	}
}

// recycleWasteToStock move all waste cards to stock, if there are are recycles available
func recycleWasteToStock(waste *Pile, stock *Pile) {
	if stock.baize.Recycles() > 0 {
		for waste.Len() > 0 {
			moveCard(waste, stock)
		}
		stock.baize.setRecycles(stock.baize.Recycles() - 1)
	}
}
