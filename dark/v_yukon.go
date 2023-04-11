package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

type Yukon struct {
	scriptBase
	extraCells int
}

func (self *Yukon) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	self.foundations = nil
	for y := 0; y < 4; y++ {
		f := self.baize.NewFoundation(newPileSlot(8, y))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

	self.cells = nil
	y := 4
	for i := 0; i < self.extraCells; i++ {
		c := self.baize.NewCell(newPileSlot(8, y))
		self.cells = append(self.cells, c)
		y += 1
	}

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 0), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownAltColor
		t.setLabel("K")
	}
}

func (self *Yukon) StartGame() {

	moveCard(self.stock, self.tableaux[0])
	var dealDown int = 1
	for x := 1; x < 7; x++ {
		for i := 0; i < dealDown; i++ {
			moveCard(self.stock, self.tableaux[x])
			if c := self.tableaux[x].peek(); c == nil {
				break
			} else {
				c.flipDown()
			}
		}
		dealDown++
		for i := 0; i < 5; i++ {
			moveCard(self.stock, self.tableaux[x])
		}
	}
}

func (*Yukon) TailMoveError([]*Card) (bool, error) {
	return true, nil
}

func (self *Yukon) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (self *Yukon) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (*Yukon) TailTapped(tail []*Card) {
	tail[0].owner().vtable.tailTapped(tail)
}

// func (*Yukon) PileTapped(*Pile) {}
