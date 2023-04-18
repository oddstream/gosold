package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"oddstream.games/gosold/cardid"
)

type FanType int

const (
	FAN_NONE FanType = iota
	FAN_DOWN
	FAN_LEFT
	FAN_RIGHT
	FAN_DOWN3
	FAN_LEFT3
	FAN_RIGHT3
)

type MoveType int

const (
	MOVE_NONE MoveType = iota
	MOVE_ANY
	MOVE_ONE
	MOVE_ONE_PLUS
	MOVE_ONE_OR_ALL
)

// pileVtabler interface for each subpile type, implements the behaviours
// specific to each subtype.
type pileVtabler interface {
	canSubtypeAppendTail([]*Card) (bool, error)
	tailTapped([]*Card)
	conformant() bool
	unsortedPairs() int
	movableTails() [][]*Card
}

type PileSlot struct {
	X, Y float32
	Deg  int
}

// Pile holds the state of the piles and cards therein.
// Pile is exported from this package because it's used to pass between light and dark.
// LIGHT should see a Pile object as immutable, hence the unexported fields and getters.
type Pile struct {
	baize                *Baize
	category             string   // needed by LIGHT when creating Pile Placeholder (switch)
	label                string   // needed by LIGHT when creating Pile Placeholder
	moveType             MoveType // needed by DARK, not visible to LIGHT
	fanType              FanType  // needed by LIGHT when fanning cards
	cards                []*Card
	vtable               pileVtabler // needed by DARK, not visible to LIGHT
	slot                 PileSlot    // needed by LIGHT when placing piles
	boundary             int         // needed by LIGHT, set by script.BuildPiles, 0 = no boundary pile
	appendFrom           string      // can only append cards from this subtype (eg Waste > Stock)
	appendCmp2, moveCmp2 dyadCmpFunc // only used by Foundation, Tableau piles
}

// Public functions, visible to LIGHT

// func (self *Pile) IsCell() bool {
// 	_, ok := self.vtable.(*Cell)
// 	return ok
// }

func (self *Pile) IsStock() bool {
	_, ok := self.vtable.(*Stock)
	return ok
}

// func (self *Pile) IsWaste() bool {
// 	_, ok := self.vtable.(*Waste)
// 	return ok
// }

func (p *Pile) Category() string {
	return p.category
}

func (p *Pile) Label() string {
	return p.label
}

func (p *Pile) Cards() []cardid.CardID {
	var ids []cardid.CardID = []cardid.CardID{}
	for _, c := range p.cards {
		ids = append(ids, c.id)
	}
	return ids
}

func (p *Pile) Boundary() int {
	return p.boundary
}

func (p *Pile) Slot() PileSlot {
	return p.slot
}

func (p *Pile) FanType() FanType {
	return p.fanType
}

// moveType is not published

// Len returns the number of cards in this pile.
// Len satisfies the sort.Interface interface.
func (self *Pile) Len() int {
	return len(self.cards)
}

func (self *Pile) Empty() bool {
	return len(self.cards) == 0
}

// Less satisfies the sort.Interface interface
func (self *Pile) Less(i, j int) bool {
	c1 := self.cards[i]
	c2 := self.cards[j]
	return c1.Suit() < c2.Suit() && c1.Ordinal() < c2.Ordinal()
}

// Swap satisfies the sort.Interface interface
func (self *Pile) Swap(i, j int) {
	self.cards[i], self.cards[j] = self.cards[j], self.cards[i]
}

// Hidden returns true if this pile is off screen
func (self *Pile) Hidden() bool {
	return self.slot.X < 0 || self.slot.Y < 0
}

// Private functions

// newPileSlot is a helper function for when slot was an image.Point
func newPileSlot(x, y int) PileSlot {
	return PileSlot{X: float32(x), Y: float32(y), Deg: 0}
}

// newHiddenPileSlot is a helper function for creating a pile (ususally stock)
// that is off-screen. The use of -5 is arbitary.
func newHiddenPileSlot() PileSlot {
	return PileSlot{X: -5, Y: -5, Deg: 0}
}

func (b *Baize) newPile(category string, slot PileSlot, fanType FanType, moveType MoveType) *Pile {
	var p *Pile = &Pile{
		baize:      b,
		category:   category,
		fanType:    fanType,
		moveType:   moveType,
		slot:       slot,
		appendCmp2: dyad.compare_Any,
		moveCmp2:   dyad.compare_Any,
	}
	b.addPile(p)
	return p
}

func (self *Pile) fill(packs, suits int) {

	self.cards = make([]*Card, 0, packs*suits*13)

	for pack := 0; pack < packs; pack++ {
		for suit := 0; suit < suits; suit++ {
			for ord := 1; ord < 14; ord++ {
				// suits are numbered NOSUIT=0, CLUB=1, DIAMOND=2, HEART=3, SPADE=4
				// (i.e. not 0..3)
				// run the suits loop backwards, so spades are used first
				// (folks expect Spider One Suit to use spades)
				var c Card = Card{id: cardid.NewCardID(pack, cardid.SPADE-suit, ord)}
				// Card will be created face up, because prone flag is not set
				// Card.pile will be nil, but set by push()
				// Card.tapTarget will be zeroed, which is what we want
				self.baize.cardMap[c.id] = &c
				self.push(&c)
			}
		}
	}

}

// func (self *Pile) addJokers(n int) {
// 	for i := 0; i < n; i++ {
// 		var c Card = Card{id: cardid.NewCardID(i, cardid.NOSUIT, 0)}
// 		self.baize.cardMap[c.id] = &c
// 		self.push(&c)
// 	}
// }

func (self *Pile) shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(self.Len(), self.Swap)
}

func (self *Pile) setLabel(label string) {
	if self.label != label {
		self.label = label
		self.baize.fnNotify(LabelEvent, self)
	}
}

// delete a *Card from this pile
func (self *Pile) delete(index int) {
	self.cards = append(self.cards[:index], self.cards[index+1:]...)
}

// extract a specific *Card from this pile
func (self *Pile) extract(pack, ordinal, suit int) *Card {
	var ID cardid.CardID = cardid.NewCardID(pack, suit, ordinal)
	for i, c := range self.cards {
		if cardid.SameCardAndPack(ID, c.id) {
			self.delete(i)
			c.flipUp()
			return c
		}
	}
	log.Printf("Could not find card %d %d in %s", suit, ordinal, self.category)
	return nil
}

func (self *Pile) extractOrdinal(ordinal int) *Card {
	for i, c := range self.cards {
		if c.Ordinal() == ordinal {
			self.delete(i)
			c.flipUp()
			return c
		}
	}
	log.Printf("Could not find card %d %s", ordinal, self.category)
	return nil
}

// peek topmost *Card of this Pile (a stack)
func (self *Pile) peek() *Card {
	if len(self.cards) == 0 {
		return nil
	}
	return self.cards[len(self.cards)-1]
}

// pop a *Card off the end of this Pile (a stack)
func (self *Pile) pop() *Card {
	if len(self.cards) == 0 {
		return nil
	}
	c := self.cards[len(self.cards)-1]
	self.cards = self.cards[:len(self.cards)-1]
	c.flipUp()
	c.setOwner(nil)
	return c
}

// push a *Card onto the end of this Pile (a stack)
func (self *Pile) push(c *Card) {
	self.cards = append(self.cards, c)
	if self.IsStock() {
		c.flipDown()
	}
	c.setOwner(self)
}

// prev returns *Card that is before specified card in the pile
// or nil if there is no card before it
func (self *Pile) prev(c *Card) *Card {
	for i, d := range self.cards {
		if c == d && i > 0 {
			return self.cards[i-1]
		}
	}
	return nil
}

func (self *Pile) flipUpExposedCard() {
	if !self.IsStock() {
		if c := self.peek(); c != nil {
			c.flipUp()
		}
	}
}

func (self *Pile) reverseCards() {
	for i, j := 0, len(self.cards)-1; i < j; i, j = i+1, j-1 {
		self.cards[i], self.cards[j] = self.cards[j], self.cards[i]
	}
}

// buryCards moves cards with the specified ordinal to the beginning of the pile
func (self *Pile) buryCards(ordinal int) {
	tmp := make([]*Card, 0, cap(self.cards))
	for _, c := range self.cards {
		if c.Ordinal() == ordinal {
			tmp = append(tmp, c)
		}
	}
	for _, c := range self.cards {
		if c.Ordinal() != ordinal {
			tmp = append(tmp, c)
		}
	}
	self.cards = self.cards[:0]
	for i := 0; i < len(tmp); i++ {
		self.push(tmp[i])
	}
}

// canMoveTail filters out cases where a tail can be moved from a given pile type
// eg if only one card can be moved at a time
// before handling off to next level of checking (the script)
// nb we skip the pile.vtable level of checking
func (self *Pile) canMoveTail(tail []*Card) (bool, error) {
	if !self.IsStock() {
		if anyCardsProne(tail) {
			return false, errors.New("Cannot move a face down card")
		}
	}
	switch self.moveType {
	case MOVE_NONE:
		// eg Discard, Foundation
		return false, fmt.Errorf("Cannot move a card from a %s", self.category)
	case MOVE_ONE:
		// eg Cell, Reserve, Stock, Waste
		if len(tail) > 1 {
			return false, fmt.Errorf("Can only move one card from a %s", self.category)
		}
	case MOVE_ANY, MOVE_ONE_PLUS:
		// MOVE_ANY: well, that was easy
		// MOVE_ONE_PLUS: don't (yet) know destination, so we allow this as MOVE_ANY
		// and do power moves check later, in Tableau canSubtypeAppendTail
	case MOVE_ONE_OR_ALL:
		// Canfield, Toad
		if len(tail) == 1 {
			// that's okay
		} else if len(tail) == self.Len() {
			// that's okay too
		} else {
			return false, errors.New("Can only move one card, or the whole pile")
		}
	}
	return self.baize.script.TailMoveError(tail)
}

func (self *Pile) canAppendTail(tail []*Card) (bool, error) {
	if self.appendFrom != "" {
		src := tail[0].owner()
		if src.category != self.appendFrom {
			return false, fmt.Errorf("A %s cannot accept cards from a %s", self.category, src.category)
		}
	}
	return self.vtable.canSubtypeAppendTail(tail)
}

func (self *Pile) makeTail(c *Card) []*Card {
	// TODO make another version that honors moveType
	// for use when making movable tails
	// would save checking made tail afterwards
	// this version would still be used when starting a tail drag
	if c.owner() != self {
		log.Panic("Pile.MakeTail called with a card that is not of this pile")
	}
	if c == self.peek() {
		return []*Card{c}
	}
	for i, pc := range self.cards {
		if pc == c {
			return self.cards[i:]
		}
	}
	log.Panicf("Pile.MakeTail could not find [%s] in %s pile", c, self.category)
	return nil
}

func (self *Pile) defaultTailTapped(tail []*Card) {
	c := tail[0]
	if dst := c.tapTarget.dst; dst != nil {
		if len(tail) == 1 {
			moveCard(c.owner(), dst)
		} else {
			moveTail(c, dst)
		}
	}
}

func (self *Pile) singleCardMovableTails() [][]*Card {
	if self.Len() > 0 {
		var card *Card = self.peek()
		var tail []*Card = []*Card{card}
		return [][]*Card{tail}
	}
	return nil
}
