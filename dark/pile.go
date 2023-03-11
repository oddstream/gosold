package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"fmt"
	"image"
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
// Made public for now, but that may change.
type pileVtabler interface {
	CanAcceptTail([]*Card) (bool, error)
	TailTapped([]*Card)
	Conformant() bool
	unsortedPairs() int
	MovableTails() []*movableTail
}

// movableTail is used for collecting tap destinations
type movableTail struct {
	dst  *Pile
	tail []*Card
}

// Pile holds the state of the piles and cards therein.
// Pile is exported from this package because it's used to pass between light and dark.
// LIGHT should see a Pile object as immutable, hence the unexported fields and getters.
type Pile struct {
	baize    *Baize
	category string   // needed by LIGHT when creating Pile Placeholder (switch)
	label    string   // needed by LIGHT when creating Pile Placeholder
	moveType MoveType // needed by DARK, not visible to LIGHT
	fanType  FanType  // needed by LIGHT when fanning cards
	cards    []*Card
	vtable   pileVtabler // needed by DARK, not visible to LIGHT
	slot     image.Point // needed by LIGHT when placing piles
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

func (p *Pile) Category() string {
	return p.category
}

func (p *Pile) Label() string {
	return p.label
}

func (p *Pile) Cards() []*Card {
	return p.cards
}

func (p *Pile) Slot() image.Point {
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

func (b *Baize) newPile(category string, slot image.Point, fanType FanType, moveType MoveType) *Pile {
	var p *Pile = &Pile{
		baize:    b,
		category: category,
		fanType:  fanType,
		moveType: moveType,
		slot:     slot,
	}
	b.addPile(p)
	return p
}

func (self *Pile) setLabel(label string) {
	if self.label != label {
		self.label = label
		self.baize.fnNotify(LabelEvent, self)
	}
}

func (self *Pile) fill(packs, suits int) int {
	var count int = packs * suits * 13

	self.cards = make([]*Card, 0, count)

	for pack := 0; pack < packs; pack++ {
		for suit := 0; suit < suits; suit++ {
			for ord := 1; ord < 14; ord++ {
				// suits are numbered NOSUIT=0, CLUB=1, DIAMOND=2, HEART=3, SPADE=4
				// (i.e. not 0..3)
				// run the suits loop backwards, so spades are used first
				// (folks expect Spider One Suit to use spades)
				var c Card = newCard(pack, cardid.SPADE-suit, ord)
				self.push(&c)
			}
		}
	}

	return count
}

func (self *Pile) shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(self.Len(), self.Swap)
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

// peek topmost Card of this Pile (a stack)
func (self *Pile) peek() *Card {
	if len(self.cards) == 0 {
		return nil
	}
	return self.cards[len(self.cards)-1]
}

// pop a Card off the end of this Pile (a stack)
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

// push a Card onto the end of this Pile (a stack)
func (self *Pile) push(c *Card) {
	self.cards = append(self.cards, c)
	if self.IsStock() {
		c.flipDown()
	}
	c.setOwner(self)
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
	case MOVE_ANY:
		// well, that was easy
	case MOVE_ONE:
		// eg Cell, Reserve, Stock, Waste
		if len(tail) > 1 {
			return false, fmt.Errorf("Can only move one card from a %s", self.category)
		}
	case MOVE_ONE_PLUS:
		// don't (yet) know destination, so we allow this as MOVE_ANY
		// and do power moves check later, in Tableau CanAcceptTail
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
	return true, nil
}

func (self *Pile) makeTail(c *Card) []*Card {
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
	card := tail[0]
	if card.tapDestination != nil {
		if len(tail) == 1 {
			moveCard(card.owner(), card.tapDestination)
		} else {
			moveTail(card, card.tapDestination)
		}
	}
}
