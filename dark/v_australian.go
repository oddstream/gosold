package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

type Australian struct {
	scriptBase
}

func (self *Australian) BuildPiles() {
	self.stock = self.baize.NewStock(newPileSlot(0, 0))
	self.waste = self.baize.NewWaste(newPileSlot(1, 0), FAN_RIGHT3)

	self.foundations = nil
	for x := 4; x < 8; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 8; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownSuit
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

func (self *Australian) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*Australian) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (self *Australian) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		c := pile.pop()
		self.waste.push(c)
	} else {
		pile.vtable.tailTapped(tail)
	}
}

// func (*Australian) PileTapped(*Pile) {}
