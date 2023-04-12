package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

type UncleSam struct {
	scriptBase
	u_foundations, s_foundations []*Pile
}

var u_coords []PileSlot = []PileSlot{
	{2, 1, 0},
	{7, 1, 0},
	{2, 2, 0},
	{3, 2, 0},
	{4, 2, 0},
	{5, 2, 0},
	{6, 2, 0},
	{7, 2, 0},
}

//                     1 1 1 1 1
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4
//
// 1   6         6   7   7 7 7
// 2   6 6 6 6 6 6   7 7 7   7
// 3

var s_coords []PileSlot = []PileSlot{
	{9, 1, 0},
	{11, 1, 0},
	{12, 1, 0},
	{13, 1, 0},
	{9, 2, 0},
	{10, 2, 0},
	{11, 2, 0},
	{13, 2, 0},
}

func (self *UncleSam) BuildPiles() {

	self.stock = self.baize.NewStock(PileSlot{0, 0, 0})
	self.waste = self.baize.NewWaste(PileSlot{0, 1, 0}, FAN_NONE)

	for _, slot := range u_coords {
		f := self.baize.NewFoundation(slot)
		self.foundations = append(self.foundations, f)
		self.u_foundations = append(self.u_foundations, f)
		f.appendCmp2 = dyad.compare_DownSuit
		f.setLabel("6")
	}

	for _, slot := range s_coords {
		f := self.baize.NewFoundation(slot)
		self.foundations = append(self.foundations, f)
		self.s_foundations = append(self.s_foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("7")
	}

	for x := 1; x < 14; x++ {
		r := self.baize.NewReserve(newPileSlot(x, 4), FAN_NONE)
		self.reserves = append(self.reserves, r)
	}
}

func (self *UncleSam) StartGame() {

	for _, f := range self.u_foundations {
		f.push(self.stock.extractOrdinal(6))
	}

	for _, f := range self.s_foundations {
		f.push(self.stock.extractOrdinal(7))
	}

	for _, r := range self.reserves {
		moveCard(self.stock, r)
	}

	// "which may be taken up and dealt ONCE after the pack in hand is exhausted"
	self.baize.setRecycles(1)
}

func (self *UncleSam) AfterMove() {
	for _, r := range self.reserves {
		if r.Empty() {
			moveCard(self.stock, r)
		}
	}
}

func (*UncleSam) TailMoveError(tail []*Card) (bool, error) {
	// attempt to move more than one card will be caught before this
	return true, nil
}

func (self *UncleSam) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (self *UncleSam) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (self *UncleSam) TailTapped(tail []*Card) {
	pile := tail[0].owner()
	switch pile.vtable.(type) {
	case *Stock:
		// move all reserve cards to waste
		for _, r := range self.reserves {
			if !r.Empty() {
				c := moveCard(r, self.waste)
				c.flipDown()
			}
		}
		// redeal reserve from stock
		// (will be done by AfterMove)
	case *Reserve:
		pile.vtable.tailTapped(tail)
	}
}

func (self *UncleSam) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.waste, self.stock)
	}
}

func (self *UncleSam) Complete() bool {
	// if the game succeeds, the states will have all been used up
	// building up the U.S., which will then show only aces and kings
	//
	// u_foundations will contain 6 5 4 3 2 1 x 8 = 6 x 8 = 48 cards
	// s_foundations will contain 7 8 9 10 J Q K x 8 = 7 x 8 = 56 cards
	// 48 + 56 = 104 = 52 x 2
	// nb there are 16 foundations, so there will only be 50% foundation occupancy
	for _, f := range self.u_foundations {
		if f.peek().Ordinal() != 1 {
			return false
		}
	}
	for _, f := range self.s_foundations {
		if f.peek().Ordinal() != 13 {
			return false
		}
	}
	return true
}
