package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"
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
	self.stock = self.baize.NewStock(image.Point{0, 0})
	self.waste = self.baize.NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	self.foundations = []*Pile{}
	for _, x := range self.founds {
		f := self.baize.NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.setLabel("A")
	}

	self.tableaux = []*Pile{}
	for _, x := range self.tabs {
		t := self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		t.setLabel("K")
		self.tableaux = append(self.tableaux, t)
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
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := tailConformant(tail, cardPair.compare_DownAltColor)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (*Klondike) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return compare_Empty(dst, tail[0])
		} else {
			return cardPair{dst.peek(), tail[0]}.compare_UpSuit()
		}
	case *Tableau:
		if dst.Empty() {
			return compare_Empty(dst, tail[0])
		} else {
			return cardPair{dst.peek(), tail[0]}.compare_DownAltColor()
		}
	}
	return true, nil
}

func (*Klondike) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.compare_DownAltColor)
}

func (self *Klondike) TailTapped(tail []*Card, nTarget int) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		for i := 0; i < self.draw; i++ {
			moveCard(self.stock, self.waste)
		}
	} else {
		pile.vtable.TailTapped(tail, nTarget)
	}
}

func (self *Klondike) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.waste, self.stock)
	}
}
