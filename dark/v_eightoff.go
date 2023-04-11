package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

type EightOff struct {
	scriptBase
}

func (self *EightOff) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	self.cells = nil
	for x := 0; x < 8; x++ {
		self.cells = append(self.cells, self.baize.NewCell(newPileSlot(x, 0)))
	}

	self.foundations = nil
	for y := 0; y < 4; y++ {
		f := self.baize.NewFoundation(newPileSlot(9, y))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 8; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ONE_PLUS)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownSuit
		t.moveCmp2 = dyad.compare_DownSuit
		t.setLabel("K")
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
	return tailConformant(tail, pile.moveCmp2)
}

func (self *EightOff) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*EightOff) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (*EightOff) TailTapped(tail []*Card) {
	tail[0].owner().vtable.tailTapped(tail)
}

// func (*EightOff) PileTapped(*Pile) {}
