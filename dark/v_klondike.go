package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"log"
)

type Klondike struct {
	scriptBase
	founds, tabs   []int
	draw, recycles int
	thoughtful     bool
}

func (self *Klondike) BuildPiles() {
	if len(self.founds) == 0 {
		self.founds = []int{3, 4, 5, 6}
	}
	if len(self.tabs) == 0 {
		self.tabs = []int{0, 1, 2, 3, 4, 5, 6}
	}
	if self.draw == 0 {
		self.draw = 1
	}
	self.stock = self.baize.NewStock(newPileSlot(0, 0))
	self.waste = self.baize.NewWaste(newPileSlot(1, 0), FAN_RIGHT3)

	self.foundations = []*Pile{}
	for _, x := range self.founds {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = cardPair.compare_UpSuit
		f.setLabel("A")
	}

	self.tableaux = []*Pile{}
	for _, x := range self.tabs {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = cardPair.compare_DownAltColor
		t.moveCmp2 = cardPair.compare_DownAltColor
		t.setLabel("K")
	}
}

func (self *Klondike) StartGame() {
	var dealDown int = 0
	for _, pile := range self.tableaux {
		for i := 0; i < dealDown; i++ {
			card := moveCard(self.stock, pile)
			if card == nil {
				log.Print("No card")
				break
			}
			if !self.thoughtful {
				card.flipDown()
			}
		}
		dealDown++
		moveCard(self.stock, pile)
	}
	self.baize.setRecycles(self.recycles)
	for i := 0; i < self.draw; i++ {
		moveCard(self.stock, self.waste)
	}
}

func (self *Klondike) AfterMove() {
	if self.waste.Len() == 0 && self.stock.Len() != 0 {
		for i := 0; i < self.draw; i++ {
			moveCard(self.stock, self.waste)
		}
	}
}

func (*Klondike) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *Klondike) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*Klondike) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(cardPair{c1, c2})
	// switch pile.vtable.(type) {
	// case *Foundation:
	// 	return cardPair{c1, c2}.compare_UpSuit()
	// case *Tableau:
	// 	return cardPair{c1, c2}.compare_DownAltColor()
	// }
	// return true, nil
}

func (self *Klondike) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		for i := 0; i < self.draw; i++ {
			moveCard(self.stock, self.waste)
		}
	} else {
		pile.vtable.tailTapped(tail)
	}
}

func (self *Klondike) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.waste, self.stock)
	}
}
