package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I damn well like, thank you

type Blockade struct {
	scriptBase
}

func (self *Blockade) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	self.foundations = nil
	for x := 4; x < 12; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.setLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 12; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
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
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := tailConformant(tail, cardPair.compare_DownSuit)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (self *Blockade) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*Blockade) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	switch pile.vtable.(type) {
	case *Foundation:
		return cardPair{c1, c2}.compare_UpSuit()
	case *Tableau:
		return cardPair{c1, c2}.compare_DownSuit()
	}
	return true, nil
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
