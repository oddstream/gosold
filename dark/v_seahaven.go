package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

type Seahaven struct {
	scriptBase
}

func (self *Seahaven) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	for x := 0; x < 4; x++ {
		self.cells = append(self.cells, self.baize.NewCell(newPileSlot(x, 0)))
	}

	for x := 6; x < 10; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

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

// default TailMoveError

// default TailAppendError

// default TailTapped

// func (*Seahaven) PileTapped(*Pile) {}
