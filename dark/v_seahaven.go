package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

type Seahaven struct {
	scriptBase
}

func (self *Seahaven) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	self.cells = nil
	for x := 0; x < 4; x++ {
		self.cells = append(self.cells, self.baize.NewCell(newPileSlot(x, 0)))
	}

	self.foundations = nil
	for x := 6; x < 10; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 10; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ONE_PLUS)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownSuit
		t.moveCmp2 = dyad.compare_DownSuit
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

func (*Seahaven) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *Seahaven) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*Seahaven) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (*Seahaven) TailTapped(tail []*Card) {
	tail[0].owner().vtable.tailTapped(tail)
}

// func (*Seahaven) PileTapped(*Pile) {}
