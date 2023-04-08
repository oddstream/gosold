package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"
)

type MrsMop struct {
	scriptBase
	easy bool
}

func (self *MrsMop) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{-5, -5})

	self.discards = []*Pile{}
	for x := 0; x < 4; x++ {
		d := self.baize.NewDiscard(image.Point{x, 0}, FAN_NONE)
		self.discards = append(self.discards, d)
		d = self.baize.NewDiscard(image.Point{x + 9, 0}, FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 13; x++ {
		t := self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
	}

	self.cells = []*Pile{}
	if self.easy {
		for x := 5; x < 8; x++ {
			t := self.baize.NewCell(image.Point{x, 0})
			self.cells = append(self.cells, t)
		}
	}
}

func (self *MrsMop) StartGame() {
	// 13 piles of 8 cards each
	for _, pile := range self.tableaux {
		for i := 0; i < 8; i++ {
			moveCard(self.stock, pile)
		}
	}
}

func (*MrsMop) TailMoveError(tail []*Card) (bool, error) {
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

func (*MrsMop) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Discard:
		// Discard.canAcceptTail() has already checked
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

func (*MrsMop) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	switch pile.vtable.(type) {
	case *Discard:
		return cardPair{c1, c2}.compare_DownSuit()
	case *Tableau:
		return cardPair{c1, c2}.compare_Down()
	}
	return true, nil
}

func (*MrsMop) TailTapped(tail []*Card) {
	tail[0].owner().vtable.tailTapped(tail)
}

// func (*MrsMop) PileTapped(*Pile) {}

func (self *MrsMop) Complete() bool {
	return self.SpiderComplete()
}
