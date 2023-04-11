package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"oddstream.games/gosold/cardid"
)

type Freecell struct {
	scriptBase
	tabCompareFunc dyadCmpFunc
	blind, easy    bool
}

func (self *Freecell) BuildPiles() {

	if self.cardColors == 0 {
		self.cardColors = 2
	}
	if self.tabCompareFunc == nil {
		self.tabCompareFunc = dyad.compare_DownAltColor
	}

	self.stock = self.baize.NewStock(newHiddenPileSlot())

	self.cells = []*Pile{}
	for x := 0; x < 4; x++ {
		c := self.baize.NewCell(newPileSlot(x, 0))
		self.cells = append(self.cells, c)
	}

	self.foundations = []*Pile{}
	for x := 4; x < 8; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 8; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ONE_PLUS)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = self.tabCompareFunc
		t.moveCmp2 = self.tabCompareFunc
	}
}

func (self *Freecell) StartGame() {
	if self.easy {
		self.foundations[0].push(self.stock.extract(0, 1, cardid.CLUB))
		self.foundations[1].push(self.stock.extract(0, 1, cardid.DIAMOND))
		self.foundations[2].push(self.stock.extract(0, 1, cardid.HEART))
		self.foundations[3].push(self.stock.extract(0, 1, cardid.SPADE))
		for _, t := range self.tableaux {
			for i := 0; i < 6; i++ {
				moveCard(self.stock, t)
			}
		}
	} else {
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
	if self.blind {
		for _, t := range self.tableaux {
			topCard := t.peek()
			for _, card := range t.cards {
				if card != topCard {
					card.flipDown()
				}
			}
		}
	}
}

func (self *Freecell) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *Freecell) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (self *Freecell) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (*Freecell) TailTapped(tail []*Card) {
	tail[0].owner().vtable.tailTapped(tail)
}

// func (*Freecell) PileTapped(*Pile) {}
