package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

type Whitehead struct {
	scriptBase
}

func (self *Whitehead) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))
	self.wastes = append(self.wastes, self.baize.NewWaste(newPileSlot(1, 0), FAN_RIGHT3))

	for x := 3; x < 7; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownAltColor
		t.moveCmp2 = dyad.compare_DownAltColor
	}
}

func (self *Whitehead) StartGame() {
	var deal = 1
	for _, pile := range self.tableaux {
		for i := 0; i < deal; i++ {
			moveCard(self.stock, pile)
		}
		deal++
	}
	self.populateWasteFromStock(1)
	self.baize.setRecycles(0)
}

func (self *Whitehead) AfterMove() {
	self.populateWasteFromStock(1)
}

// default TailMoveError

// default TailAppendError

// default TailTapped

func (self *Whitehead) PileTapped(*Pile) {
	// https://politaire.com/help/whitehead
	// Only one pass through the Stock is permitted
}
