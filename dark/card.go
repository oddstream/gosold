package dark

import (
	"oddstream.games/gosold/cardid"
)

type tapTarget struct {
	dst     *Pile // where this card will go if tapped
	weight  int16 // 0...4; how fruity this move is
	percent int   // what the Baize would be if this tap is tapped
}

// Card holds the state of the cards.
// Card is exported from this package because it's used to pass between light and dark.
// LIGHT should see a Card object as immutable, hence the unexported fields and getters.
type Card struct {
	id   cardid.CardID
	pile *Pile // the Pile this Card is currently in
	tapTarget
}

// Public functions, visible outside DARK

func (c *Card) String() string {
	return c.id.String()
}

// func (c *Card) ID() cardid.CardID {
// 	return c.id
// }

func (c *Card) Pack() int {
	return c.id.Pack()
}

func (c *Card) Suit() int {
	return c.id.Suit()
}

func (c *Card) Ordinal() int {
	return c.id.Ordinal()
}

func (c *Card) Prone() bool {
	return c.id.Prone()
}

func (c *Card) Black() bool {
	return c.id.Black()
}

func (c *Card) TapWeight() int16 {
	return c.tapTarget.weight
}

func (c *Card) TapPercent() int {
	return c.tapTarget.percent
}

// func (c *Card) SetProne(prone bool) {
// 	c.id = c.id.SetProne(prone)
// }

// Private functions, only visible inside DARK

func (c *Card) owner() *Pile {
	return c.pile
}

func (c *Card) setOwner(p *Pile) {
	c.pile = p
}

func (c *Card) setProne(prone bool) {
	c.id = c.id.SetProne(prone)
}

func (c *Card) flipUp() {
	if c.Prone() {
		c.setProne(false)
	}
}

func (c *Card) flipDown() {
	if !c.Prone() {
		c.setProne(true)
	}
}
