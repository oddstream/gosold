package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

type MrsMop struct {
	scriptBase
	easy bool
}

func (self *MrsMop) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	self.discards = []*Pile{}
	for x := 0; x < 4; x++ {
		d := self.baize.NewDiscard(newPileSlot(x, 0), FAN_NONE)
		self.discards = append(self.discards, d)
		d = self.baize.NewDiscard(newPileSlot(x+9, 0), FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 13; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownSuit
		t.moveCmp2 = dyad.compare_DownSuit
	}

	self.cells = []*Pile{}
	if self.easy {
		for x := 5; x < 8; x++ {
			c := self.baize.NewCell(newPileSlot(x, 0))
			self.cells = append(self.cells, c)
		}
	}
}

func (self *MrsMop) StartGame() {
	// 13 piles of 8 cards each
	for _, pile := range self.tableaux {
		for i := 0; i < 8; i++ {
			moveCard(self.stock, pile)
		}
	}
}

func (*MrsMop) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *MrsMop) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*MrsMop) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (*MrsMop) TailTapped(tail []*Card) {
	tail[0].owner().vtable.tailTapped(tail)
}

// func (*MrsMop) PileTapped(*Pile) {}

func (self *MrsMop) Complete() bool {
	return self.SpiderComplete()
}
