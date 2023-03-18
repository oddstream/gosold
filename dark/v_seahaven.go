package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"
)

type Seahaven struct {
	scriptBase
}

func (self *Seahaven) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{-5, -5})

	self.cells = nil
	for x := 0; x < 4; x++ {
		self.cells = append(self.cells, self.baize.NewCell(image.Point{x, 0}))
	}

	self.foundations = nil
	for x := 6; x < 10; x++ {
		f := self.baize.NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.setLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 10; x++ {
		t := self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS)
		self.tableaux = append(self.tableaux, t)
		t.setLabel("K")
	}
}

func (self *Seahaven) StartGame() {
	for _, t := range self.tableaux {
		for i := 0; i < 5; i++ {
			moveCard(self.stock, t)
		}
	}
	moveCard(self.stock, self.cells[1])
	moveCard(self.stock, self.cells[2])
	self.baize.setRecycles(0)
}

func (self *Seahaven) TailMoveError(tail []*Card) (bool, error) {
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

func (self *Seahaven) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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

func (*Seahaven) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.compare_DownSuit)
}

func (*Seahaven) TailTapped(tail []*Card, nTarget int) {
	tail[0].owner().vtable.TailTapped(tail, nTarget)
}

// func (*Seahaven) PileTapped(*Pile) {}
