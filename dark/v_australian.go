package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"
)

type Australian struct {
	scriptBase
}

func (self *Australian) BuildPiles() {
	self.stock = self.baize.NewStock(image.Point{0, 0})
	self.waste = self.baize.NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	self.foundations = nil
	for x := 4; x < 8; x++ {
		f := self.baize.NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.setLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 8; x++ {
		t := self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.setLabel("K")
	}
}

func (self *Australian) StartGame() {
	for _, pile := range self.tableaux {
		for i := 0; i < 4; i++ {
			moveCard(self.stock, pile)
		}
	}
	moveCard(self.stock, self.waste)
	self.baize.setRecycles(0)
}

func (self *Australian) AfterMove() {
	if self.waste.Len() == 0 && self.stock.Len() != 0 {
		moveCard(self.stock, self.waste)
	}
}

func (*Australian) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (*Australian) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return compare_Empty(dst, tail[0])
		} else {
			return cardPair{dst.peek(), tail[0]}.compare_UpSuit()
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

func (*Australian) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.compare_DownSuit)
}

func (self *Australian) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		c := pile.pop()
		self.waste.push(c)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

// func (*Australian) PileTapped(*Pile) {}
