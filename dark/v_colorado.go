package dark

import (
	"errors"

	"oddstream.games/gosold/util"
)

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

type Colorado struct {
	scriptBase
}

func (self *Colorado) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	self.reserves = append(self.reserves, self.baize.NewReserve(newPileSlot(0, 0), FAN_NONE))

	for x := 2; x < 6; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}
	for x := 6; x < 10; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_DownSuit
		f.setLabel("K")
	}

	for x := 0; x < 10; x++ {
		for y := 1; y < 4; y += 2 {
			w := self.baize.NewWaste(newPileSlot(x, y), FAN_DOWN3)
			self.wastes = append(self.wastes, w)
			w.appendFrom = "Reserve" // default is "Stock"
		}
	}
}

func (self *Colorado) StartGame() {

	for _, w := range self.wastes {
		moveCard(self.stock, w)
	}

	for !self.stock.Empty() {
		moveCard(self.stock, self.reserves[0])
	}

	self.baize.setRecycles(0)
}

func (self *Colorado) AfterMove() {
	for _, t := range self.wastes {
		if t.Empty() {
			moveCard(self.reserves[0], t)
		}
	}
}

// func (*Colorado) TailMoveError(tail []*Card) (bool, error) {
// 	return true, nil
// }

func (self *Colorado) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail)
	}
	src := tail[0].pile
	if util.Contains(self.wastes, src) && util.Contains(self.wastes, dst) {
		return false, errors.New("Cannot move cards between Waste piles")
	}
	return dst.appendCmp2(dyad{dst.peek(), tail[0]})
}

// default TailTapped

// func (*Colorado) PileTapped(*Pile) {}
