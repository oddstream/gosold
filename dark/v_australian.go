package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

type Australian struct {
	scriptBase
}

func (self *Australian) BuildPiles() {
	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	self.wastes = append(self.wastes, self.baize.NewWaste(newPileSlot(1, 0), FAN_RIGHT3))

	for x := 4; x < 8; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

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
	self.populateWasteFromStock(1)
	self.baize.setRecycles(0)
}

func (self *Australian) AfterMove() {
	self.populateWasteFromStock(1)
}

func (*Australian) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

// default TailAppendError

// default TailTapped

// func (*Australian) PileTapped(*Pile) {}
