package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

type Whitehead struct {
	scriptBase
}

func (self *Whitehead) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))
	self.waste = self.baize.NewWaste(newPileSlot(1, 0), FAN_RIGHT3)

	self.foundations = nil
	for x := 3; x < 7; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = cardPair.compare_UpSuit
		f.setLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = cardPair.compare_DownAltColor
		t.moveCmp2 = cardPair.compare_DownAltColor
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
	self.baize.setRecycles(0)
	moveCard(self.stock, self.waste)
}

func (self *Whitehead) AfterMove() {
	if self.waste.Len() == 0 && self.stock.Len() != 0 {
		moveCard(self.stock, self.waste)
	}
}

func (*Whitehead) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *Whitehead) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (self *Whitehead) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(cardPair{c1, c2})
	// switch pile.vtable.(type) {
	// case *Foundation:
	// 	return cardPair{c1, c2}.compare_UpSuit()
	// case *Tableau:
	// 	return cardPair{c1, c2}.compare_DownColor()
	// }
	// return true, nil
}

func (self *Whitehead) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		moveCard(self.stock, self.waste)
	} else {
		pile.vtable.tailTapped(tail)
	}
}

func (self *Whitehead) PileTapped(*Pile) {
	// https://politaire.com/help/whitehead
	// Only one pass through the Stock is permitted
}
