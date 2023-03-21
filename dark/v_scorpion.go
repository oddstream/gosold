package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"
)

type Scorpion struct {
	scriptBase
}

func (self *Scorpion) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{0, 0})

	self.discards = []*Pile{}
	for x := 3; x < 7; x++ {
		d := self.baize.NewDiscard(image.Point{x, 0}, FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		t.setLabel("K")
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Scorpion) StartGame() {
	// The Tableau consists of 10 stacks with 6 cards in the first 4 stacks, with the 6th card face up,
	// and 5 cards in the remaining 6 stacks, with the 5th card face up.

	for _, tab := range self.tableaux {
		for i := 0; i < 7; i++ {
			moveCard(self.stock, tab)
		}
	}

	for i := 0; i < 4; i++ {
		tab := self.tableaux[i]
		for j := 0; j < 3; j++ {
			tab.cards[j].flipDown()
		}
	}
	self.baize.setRecycles(0)
}

func (*Scorpion) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (*Scorpion) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch dst.vtable.(type) {
	case *Discard:
		if tail[0].Ordinal() != 13 {
			return false, errors.New("Can only discard starting from a King")
		}
		ok, err := tailConformant(tail, cardPair.compare_DownSuit)
		if !ok {
			return ok, err
		}
	case *Tableau:
		if dst.Empty() {
			return compare_Empty(dst, tail[0])
		} else {
			return cardPair{dst.peek(), tail[0]}.compare_DownSuit()
		}
	}
	return true, nil
}

func (*Scorpion) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.compare_DownSuit)
}

func (self *Scorpion) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	switch pile.vtable.(type) {
	case *Stock:
		if !self.stock.Empty() {
			for _, tab := range self.tableaux {
				moveCard(self.stock, tab)
			}
		}
	default:
		tail[0].owner().vtable.TailTapped(tail)
	}
}

// func (*Scorpion) PileTapped(*Pile) {}

func (self *Scorpion) Complete() bool {
	return self.SpiderComplete()
}
