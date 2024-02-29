package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

type Freecell struct {
	scriptBase
}

func (self *Freecell) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	for x := 0; x < 4; x++ {
		c := self.baize.NewCell(newPileSlot(x, 0))
		self.cells = append(self.cells, c)
	}

	for x := 4; x < 8; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmpFunc = dyad.compare_UpSuit
		f.setLabel("A")
	}

	for x := 0; x < 8; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ONE_PLUS)
		self.tableaux = append(self.tableaux, t)
		t.appendCmpFunc = dyad.compare_DownAltColor
		t.moveCmpFunc = dyad.compare_DownAltColor
	}
}

func (self *Freecell) StartGame() {
	// 4 piles of 7 cards
	// 4 piles of 6 cards
	for i := 0; i < 4; i++ {
		t := self.tableaux[i]
		for j := 0; j < 7; j++ {
			moveCard(self.stock, t)
		}
	}
	for i := 4; i < 8; i++ {
		t := self.tableaux[i]
		for j := 0; j < 6; j++ {
			moveCard(self.stock, t)
		}
	}
}

// default TailMoveError

// default TailAppendError

// default TailTapped

// default PileTapped
