package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

type Scorpion struct {
	scriptBase
}

func (self *Scorpion) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	for x := 3; x < 7; x++ {
		d := self.baize.NewDiscard(newPileSlot(x, 0), FAN_NONE)
		self.discards = append(self.discards, d)
	}

	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownSuit
		t.setLabel("K")
	}
}

func (self *Scorpion) StartGame() {
	// The Tableau consists of 10 stacks with 6 cards in the first 4 stacks, with the 6th card face up,
	// and 5 cards in the remaining 6 stacks, with the 5th card face up.

	for _, tab := range self.tableaux {
		for i := 0; i < 7; i++ {
			moveCard(self.stock, tab)
		}
	}

	for i := 0; i < 4; i++ {
		tab := self.tableaux[i]
		for j := 0; j < 3; j++ {
			tab.cards[j].flipDown()
		}
	}
	self.baize.setRecycles(0)
}

func (*Scorpion) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (self *Scorpion) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*Scorpion) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (self *Scorpion) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	switch pile.vtable.(type) {
	case *Stock:
		if !self.stock.Empty() {
			for _, tab := range self.tableaux {
				moveCard(self.stock, tab)
			}
		}
	default:
		tail[0].owner().vtable.tailTapped(tail)
	}
}

// func (*Scorpion) PileTapped(*Pile) {}

func (self *Scorpion) Complete() bool {
	return self.SpiderComplete()
}
