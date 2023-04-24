package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

type SimpleSimon struct {
	scriptBase
}

func (self *SimpleSimon) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	for x := 3; x < 7; x++ {
		d := self.baize.NewDiscard(newPileSlot(x, 0), FAN_NONE)
		self.discards = append(self.discards, d)
	}

	for x := 0; x < 10; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_Down
		t.moveCmp2 = dyad.compare_DownSuit
	}
}

func (self *SimpleSimon) StartGame() {
	// 3 piles of 8 cards each
	for i := 0; i < 3; i++ {
		pile := self.tableaux[i]
		for j := 0; j < 8; j++ {
			moveCard(self.stock, pile)
		}
	}
	var deal int = 7
	for i := 3; i < 10; i++ {
		pile := self.tableaux[i]
		for j := 0; j < deal; j++ {
			moveCard(self.stock, pile)
		}
		deal--
	}
}

// default TailMoveError

// default TailAppendError

// default TailTapped

// func (*SimpleSimon) PileTapped(*Pile) {}

func (self *SimpleSimon) Complete() bool {
	return self.SpiderComplete()
}
