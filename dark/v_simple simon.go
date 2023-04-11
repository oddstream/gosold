package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

type SimpleSimon struct {
	scriptBase
}

func (self *SimpleSimon) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	self.discards = []*Pile{}
	for x := 3; x < 7; x++ {
		d := self.baize.NewDiscard(newPileSlot(x, 0), FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 10; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = cardPair.compare_DownSuit
		t.moveCmp2 = cardPair.compare_DownSuit
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

func (*SimpleSimon) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *SimpleSimon) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
	// switch dst.vtable.(type) {
	// case *Discard:
	// 	// Discard.canAcceptTail() has already checked
	// 	// (1) pile is empty
	// 	// (2) no prone cards in tail
	// 	// (3) tail is the length of a complete set (eg 13)
	// 	if tail[0].Ordinal() != 13 {
	// 		return false, errors.New("Can only discard starting from a King")
	// 	}
	// 	return tailConformant(tail, cardPair.compare_DownSuit)
	// case *Tableau:
	// 	if dst.Empty() {
	// 	} else {
	// 		return self.TwoCards(dst, dst.peek(), tail[0])
	// 	}
	// }
	// return true, nil
}

func (*SimpleSimon) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(cardPair{c1, c2})
	// switch pile.vtable.(type) {
	// case *Discard:
	// 	return cardPair{c1, c2}.compare_DownSuit()
	// case *Tableau:
	// 	return cardPair{c1, c2}.compare_Down()
	// }
	// return true, nil
}

func (*SimpleSimon) TailTapped(tail []*Card) {
	tail[0].owner().vtable.tailTapped(tail)
}

// func (*SimpleSimon) PileTapped(*Pile) {}

func (self *SimpleSimon) Complete() bool {
	return self.SpiderComplete()
}
