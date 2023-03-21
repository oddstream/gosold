package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"
)

type Yukon struct {
	scriptBase
	extraCells int
}

func (self *Yukon) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{-5, -5})

	self.foundations = nil
	for y := 0; y < 4; y++ {
		f := self.baize.NewFoundation(image.Point{8, y})
		self.foundations = append(self.foundations, f)
		f.setLabel("A")
	}

	self.cells = nil
	y := 4
	for i := 0; i < self.extraCells; i++ {
		c := self.baize.NewCell(image.Point{8, y})
		self.cells = append(self.cells, c)
		y += 1
	}

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(image.Point{x, 0}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.setLabel("K")
	}
}

func (self *Yukon) StartGame() {

	moveCard(self.stock, self.tableaux[0])
	var dealDown int = 1
	for x := 1; x < 7; x++ {
		for i := 0; i < dealDown; i++ {
			moveCard(self.stock, self.tableaux[x])
			if c := self.tableaux[x].peek(); c == nil {
				break
			} else {
				c.flipDown()
			}
		}
		dealDown++
		for i := 0; i < 5; i++ {
			moveCard(self.stock, self.tableaux[x])
		}
	}
}

func (*Yukon) TailMoveError([]*Card) (bool, error) {
	return true, nil
}

func (*Yukon) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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
			return cardPair{dst.peek(), tail[0]}.compare_DownAltColor()
		}
	}
	return true, nil
}

func (*Yukon) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.compare_DownAltColor)
}

func (*Yukon) TailTapped(tail []*Card) {
	tail[0].owner().vtable.TailTapped(tail)
}

// func (*Yukon) PileTapped(*Pile) {}
