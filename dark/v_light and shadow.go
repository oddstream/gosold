package dark

import (
	"errors"

	"oddstream.games/gosold/util"
)

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

type LightAndShadow struct {
	scriptBase
	auxilliaries, rivals []*Pile
}

func (self *LightAndShadow) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	self.wastes = append(self.wastes, self.baize.NewWaste(newPileSlot(0, 1), FAN_DOWN3))

	// auxilliaries
	for x := 2; x < 6; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 0), FAN_DOWN, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		self.auxilliaries = append(self.auxilliaries, t)
		t.appendCmp2 = dyad.compare_DownAltColor
		t.moveCmp2 = dyad.compare_DownAltColor
	}

	for i := 0; i < 4; i++ {
		self.tableaux[i].boundary = 2 + 4 + i
	}

	// rivals
	for x := 2; x < 6; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 2), FAN_NONE, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		self.rivals = append(self.rivals, t)
		t.appendCmp2 = dyad.compare_DownAltColor
		t.moveCmp2 = dyad.compare_DownAltColor
	}

	// foundations
	for x := 0; x < 8; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 3))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpAltColor
		f.setLabel("A")
	}

}

func (self *LightAndShadow) StartGame() {

	for _, f := range self.foundations {
		if c := self.stock.extractOrdinal(1); c != nil {
			f.push(c)
		}
	}

	for _, t := range self.tableaux {
		moveCard(self.stock, t)
	}

	self.populateWasteFromStock(1)
	self.baize.setRecycles(0)
}

func (self *LightAndShadow) AfterMove() {
	self.populateWasteFromStock(1)
}

func (*LightAndShadow) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *LightAndShadow) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	src := tail[0].pile
	if util.Contains(self.auxilliaries, dst) {
		if dst.Empty() {
			if !util.Contains(self.rivals, src) {
				return false, errors.New("Vacancies in the auxilliaries are filled from the rivals")
			}
		}
	}
	if util.Contains(self.rivals, dst) {
		if dst.Empty() {
			if src != self.Waste() {
				return false, errors.New("Vacancies in the rivals are filled from the waste")
			}
		} else {
			return false, errors.New("Rivals can only contain one card")
		}
	}
	if util.Contains(self.foundations, dst) {
		if !util.Contains(self.auxilliaries, src) {
			return false, errors.New("Foundation cards must come from the auxilliaries")
		}
	}
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])

	// TODO BUG can't tap a card in the rivals to send it to auxilliaries
}

func (*LightAndShadow) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (self *LightAndShadow) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		moveCard(self.stock, self.Waste())
	} else {
		pile.vtable.tailTapped(tail)
	}
}

func (self *LightAndShadow) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.Waste(), self.stock)
	}
}
