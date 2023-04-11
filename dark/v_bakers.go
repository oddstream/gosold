package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

type BakersDozen struct {
	scriptBase
}

func (self *BakersDozen) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 0), FAN_DOWN, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_Down
		t.setLabel("X")
	}
	for x := 0; x < 6; x++ {
		t := self.baize.NewTableau(PileSlot{float32(x) + 0.5, 3, 0}, FAN_DOWN, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_Down
		t.setLabel("X")
	}
	for x := 0; x < 6; x++ {
		// stock is pile index 0
		// tableaux are piles index 1 .. 13
		self.tableaux[x].boundary = 1 + x + 7
	}
	self.tableaux[6].boundary = 1 + 12

	self.foundations = nil
	for y := 0; y < 4; y++ {
		f := self.baize.NewFoundation(newPileSlot(9, y))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
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
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*BakersDozen) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (*BakersDozen) TailTapped(tail []*Card) {
	tail[0].owner().vtable.tailTapped(tail)
}

// func (*BakersDozen) PileTapped(*Pile) {}
