package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"oddstream.games/gosold/util"
)

type Chameleon struct {
	scriptBase
}

func (self *Chameleon) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(5, 0))

	self.wastes = append(self.wastes, self.baize.NewWaste(newPileSlot(5, 1), FAN_DOWN3))

	for x := 0; x < 4; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuitWrap
	}

	self.reserves = append(self.reserves, self.baize.NewReserve(newPileSlot(0, 1), FAN_NONE))

	for x := 1; x < 4; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownWrap
	}

}

func (self *Chameleon) StartGame() {

	for i := 0; i < 12; i++ {
		moveCard(self.stock, self.reserves[0])
	}
	card := moveCard(self.stock, self.foundations[0])
	for _, pile := range self.foundations {
		pile.setLabel(util.OrdinalToShortString(card.Ordinal()))
	}

	for _, pile := range self.tableaux {
		moveCard(self.stock, pile)
	}

	moveCard(self.stock, self.Waste())

	self.baize.setRecycles(0)
}

func (self *Chameleon) AfterMove() {
	// "fill each [tableau] space at once with the top card of the reserve,
	// after the reserve is exhausted, fill spaces from the waste pile,
	// but at this time a space may be kept open for as long as desired"
	for _, pile := range self.tableaux {
		if pile.Empty() {
			moveCard(self.reserves[0], pile)
		}
	}
	self.populateWasteFromStock(1)
}

func (self *Chameleon) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *Chameleon) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (self *Chameleon) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (self *Chameleon) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		moveCard(self.stock, self.Waste())
	} else {
		pile.vtable.tailTapped(tail)
	}
}

func (self *Chameleon) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.Waste(), self.stock)
	}
}
