package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"oddstream.games/gosold/util"
)

type Agnes struct {
	scriptBase
}

func (self *Agnes) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	self.foundations = nil
	for x := 3; x < 7; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuitWrap
	}

	self.reserves = nil
	for x := 0; x < 7; x++ {
		r := self.baize.NewReserve(newPileSlot(x, 1), FAN_NONE)
		self.reserves = append(self.reserves, r)
	}

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 2), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownAltColorWrap
		t.moveCmp2 = dyad.compare_DownAltColorWrap
	}
}

func (self *Agnes) StartGame() {

	for _, pile := range self.reserves {
		moveCard(self.stock, pile)
	}

	var dealDown int = 0
	for _, pile := range self.tableaux {
		for i := 0; i < dealDown; i++ {
			card := moveCard(self.stock, pile)
			card.flipDown()
		}
		dealDown++
		moveCard(self.stock, pile)
	}

	c := moveCard(self.stock, self.foundations[0])
	ord := c.Ordinal()
	for _, pile := range self.foundations {
		pile.setLabel(util.OrdinalToShortString(ord))
	}
	ord -= 1
	if ord == 0 {
		ord = 13
	}
	for _, pile := range self.tableaux {
		pile.setLabel(util.OrdinalToShortString(ord))
	}
}

func (self *Agnes) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *Agnes) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*Agnes) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (self *Agnes) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		for _, pile := range self.reserves {
			moveCard(self.stock, pile)
		}
	} else {
		pile.vtable.tailTapped(tail)
	}
}

// func (*Agnes) PileTapped(*Pile) {}
