package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"
)

type EightOff struct {
	scriptBase
}

func (self *EightOff) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{-5, -5})

	self.cells = nil
	for x := 0; x < 8; x++ {
		self.cells = append(self.cells, self.baize.NewCell(image.Point{x, 0}))
	}

	self.foundations = nil
	for y := 0; y < 4; y++ {
		pile := self.baize.NewFoundation(image.Point{9, y})
		self.foundations = append(self.foundations, pile)
		pile.setLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 8; x++ {
		pile := self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS)
		self.tableaux = append(self.tableaux, pile)
		pile.setLabel("K")
	}
}

func (self *EightOff) StartGame() {
	for i := 0; i < 4; i++ {
		moveCard(self.stock, self.cells[i])
	}
	for _, pile := range self.tableaux {
		for i := 0; i < 6; i++ {
			moveCard(self.stock, pile)
		}
	}
}

func (*EightOff) TailMoveError(tail []*Card) (bool, error) {
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

func (*EightOff) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
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

func (*EightOff) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.compare_DownSuit)
}

func (*EightOff) TailTapped(tail []*Card, nTarget int) {
	tail[0].owner().vtable.TailTapped(tail, nTarget)
}

// func (*EightOff) PileTapped(*Pile) {}
