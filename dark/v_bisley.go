package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"oddstream.games/gosold/cardid"
)

type Bisley struct {
	scriptBase
}

func (self *Bisley) BuildPiles() {

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	self.foundations = nil

	for x := 0; x < 4; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = cardPair.compare_DownSuit
		f.setLabel("K")
	}

	for x := 0; x < 4; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 1))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = cardPair.compare_UpSuit
		f.setLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 13; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 2), FAN_DOWN, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = cardPair.compare_UpOrDownSuit
		t.setLabel("X")
	}
}

func (self *Bisley) StartGame() {

	self.foundations[4].push(self.stock.extract(0, 1, cardid.CLUB))
	self.foundations[5].push(self.stock.extract(0, 1, cardid.DIAMOND))
	self.foundations[6].push(self.stock.extract(0, 1, cardid.HEART))
	self.foundations[7].push(self.stock.extract(0, 1, cardid.SPADE))

	// the first 4 tableaux have 3 cards
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			moveCard(self.stock, self.tableaux[i])
		}
	}
	// the next 9 tableaux have 4 cards
	for i := 4; i < 13; i++ {
		for j := 0; j < 4; j++ {
			moveCard(self.stock, self.tableaux[i])
		}
	}

	self.baize.setRecycles(0)
}

func (*Bisley) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (self *Bisley) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*Bisley) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(cardPair{c1, c2})
	// switch pile.vtable.(type) {
	// case *Foundation:
	// 	if pile.Label() == "A" {
	// 		return cardPair{c1, c2}.compare_UpSuit()
	// 	} else {
	// 		return cardPair{c1, c2}.compare_DownSuit()
	// 	}
	// case *Tableau:
	// 	return cardPair{c1, c2}.chainCall(cardPair.compare_UpOrDown, cardPair.compare_Suit)
	// }
	// return true, nil
}

func (self *Bisley) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		moveCard(self.stock, self.waste)
	} else {
		pile.vtable.tailTapped(tail)
	}
}

// func (*Bisley) PileTapped(*Pile) {}
