package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
)

type BakersDozen struct {
	scriptBase
}

func (self *BakersDozen) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{-5, -5})

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(image.Point{x, 0}, FAN_DOWN, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		t.setLabel("X")
	}
	for x := 0; x < 6; x++ {
		t := self.baize.NewTableau(image.Point{x, 3}, FAN_DOWN, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		t.setLabel("X")
	}

	self.foundations = nil
	for y := 0; y < 4; y++ {
		f := self.baize.NewFoundation(image.Point{9, y})
		self.foundations = append(self.foundations, f)
		f.setLabel("A")
	}
}

func (self *BakersDozen) StartGame() {

	for _, tab := range self.tableaux {
		for x := 0; x < 4; x++ {
			moveCard(self.stock, tab)
		}
		// demote kings
		tab.buryCards(13)
	}
}

func (*BakersDozen) TailMoveError(tail []*Card) (bool, error) {
	// attempt to move more than one card will be caught before this
	return true, nil
}

func (self *BakersDozen) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		switch dst.vtable.(type) {
		case *Foundation:
			return compare_Empty(dst, tail[0])
		case *Tableau:
			return false, errors.New("Cannot move a card to an empty Tableau")
		}
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*BakersDozen) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	switch pile.vtable.(type) {
	case *Foundation:
		return cardPair{c1, c2}.compare_UpSuit()
	case *Tableau:
		return cardPair{c1, c2}.compare_Down()
	}
	return true, nil
}

func (*BakersDozen) TailTapped(tail []*Card) {
	tail[0].owner().vtable.TailTapped(tail)
}

// func (*BakersDozen) PileTapped(*Pile) {}
