package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I damn well like, thank you

type UskPileInfo struct {
	x int
	n int
}

type Usk struct {
	scriptBase
	tableauLabel string
	layout       []UskPileInfo
}

func (self *Usk) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	self.layout = []UskPileInfo{
		{x: 0, n: 8},
		{x: 1, n: 8},
		{x: 2, n: 8},
		{x: 3, n: 7},
		{x: 4, n: 6},
		{x: 5, n: 5},
		{x: 6, n: 4},
		{x: 7, n: 3},
		{x: 8, n: 2},
		{x: 9, n: 1},
	}

	for x := 6; x < 10; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

	for _, li := range self.layout {
		t := self.baize.NewTableau(newPileSlot(li.x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownAltColor
		t.moveCmp2 = dyad.compare_DownAltColor
		t.setLabel(self.tableauLabel)
	}
}

func (self *Usk) dealCards() {
	for i, li := range self.layout {
		var t *Pile = self.tableaux[i]
		for n := 0; n < li.n; n++ {
			moveCard(self.stock, t)
		}
	}
}

func (self *Usk) StartGame() {
	self.dealCards()
	self.baize.setRecycles(1)
	// if self.tableauLabel == "" {
	// TheGame.UI.ToastInfo("Relaxed version - any card may be placed in an empty tableaux pile")
	// }
}

func (*Usk) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *Usk) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (self *Usk) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (*Usk) TailTapped(tail []*Card) {
	tail[0].owner().vtable.tailTapped(tail)
}

func (self *Usk) PileTapped(pile *Pile) {
	if pile != self.stock {
		return
	}
	if self.baize.Recycles() == 0 {
		// TheGame.UI.ToastError("No more recycles")
		return
	}
	/*
		The redeal procedure begins by picking up all cards on the tableau.
		The cards from the tableau are collected, one column at a time,
		starting with the left-most column,
		picking up the cards in each column in bottom to top order.
		Then, without shuffling, the cards are dealt out again,
		starting with the first card picked up.
		Deal the tableau in the same arrangement as it was originally dealt,
		one row at a time, starting with the bottom-most row,
		dealing the cards in each row in left to right order.
	*/
	// collect cards to stock
	for _, t := range self.tableaux {
		// moveCards keeps the order of the cards
		if !t.Empty() {
			moveTail(t.cards[0], self.stock)
		}
	}
	// reverse order so we can pop
	self.stock.reverseCards()
	// redeal cards
	self.dealCards()
	self.baize.setRecycles(0)
}
