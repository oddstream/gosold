package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"
)

type SimpleSimon struct {
	scriptBase
}

func (self *SimpleSimon) BuildPiles() {

	self.stock = NewStock(image.Point{-5, -5}, FAN_NONE, 1, 4, nil, 0)

	self.discards = []*Pile{}
	for x := 3; x < 7; x++ {
		d := NewDiscard(image.Point{x, 0}, FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 10; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *SimpleSimon) StartGame() {
	// 3 piles of 8 cards each
	for i := 0; i < 3; i++ {
		pile := self.tableaux[i]
		for j := 0; j < 8; j++ {
			moveCard(self.stock, pile)
		}
	}
	var deal int = 7
	for i := 3; i < 10; i++ {
		pile := self.tableaux[i]
		for j := 0; j < deal; j++ {
			moveCard(self.stock, pile)
		}
		deal--
	}
}

func (*SimpleSimon) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := tailConformant(tail, cardPair.compare_DownSuit)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (*SimpleSimon) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Discard:
		// Discard.CanAcceptTail() has already checked
		// (1) pile is empty
		// (2) no prone cards in tail
		// (3) tail is the length of a complete set (eg 13)
		if tail[0].Ordinal() != 13 {
			return false, errors.New("Can only discard starting from a King")
		}
		ok, err := tailConformant(tail, cardPair.compare_DownSuit)
		if !ok {
			return ok, err
		}
	case *Tableau:
		if dst.Empty() {
		} else {
			return cardPair{dst.peek(), tail[0]}.compare_Down()
		}
	}
	return true, nil
}

func (*SimpleSimon) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.compare_DownSuit)
}

func (*SimpleSimon) TailTapped(tail []*Card) {
	tail[0].owner().vtable.TailTapped(tail)
}

// func (*SimpleSimon) PileTapped(*Pile) {}

func (self *SimpleSimon) Complete() bool {
	return self.SpiderComplete()
}
