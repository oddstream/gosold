package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I damn well like, thank you

type Blockade struct {
	scriptBase
}

func (self *Blockade) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	for x := 4; x < 12; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

	for x := 0; x < 12; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownSuit
		t.moveCmp2 = dyad.compare_DownSuit
	}
}

func (self *Blockade) StartGame() {
	for _, pile := range self.tableaux {
		moveCard(self.stock, pile)
	}
	self.baize.setRecycles(0)
}

func (self *Blockade) AfterMove() {
	// An empty pile will be filled up immediately by a card from the stock.
	for _, pile := range self.tableaux {
		if pile.Empty() {
			moveCard(self.stock, pile)
		}
	}
}

func (*Blockade) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *Blockade) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*Blockade) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (self *Blockade) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock {
		for _, tab := range self.tableaux {
			moveCard(self.stock, tab)
		}
	} else {
		pile.vtable.tailTapped(tail)
	}
}

// func (*Blockade) PileTapped(*Pile) {}
