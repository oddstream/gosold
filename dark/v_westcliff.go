package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"oddstream.games/gosold/cardid"
)

type Westcliff struct {
	scriptBase
	variant string
}

func (self *Westcliff) BuildPiles() {
	self.stock = self.baize.NewStock(newPileSlot(0, 0))
	switch self.variant {
	case "Classic":
		self.waste = self.baize.NewWaste(newPileSlot(1, 0), FAN_RIGHT3)
		self.foundations = []*Pile{}
		for x := 3; x < 7; x++ {
			f := self.baize.NewFoundation(newPileSlot(x, 0))
			self.foundations = append(self.foundations, f)
			f.appendCmp2 = dyad.compare_UpSuit
			f.setLabel("A")
		}
		self.tableaux = []*Pile{}
		for x := 0; x < 7; x++ {
			t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
			self.tableaux = append(self.tableaux, t)
			t.appendCmp2 = dyad.compare_DownAltColor
			t.moveCmp2 = dyad.compare_DownAltColor
		}
	case "American":
		self.waste = self.baize.NewWaste(newPileSlot(1, 0), FAN_RIGHT3)
		self.foundations = []*Pile{}
		for x := 6; x < 10; x++ {
			f := self.baize.NewFoundation(newPileSlot(x, 0))
			self.foundations = append(self.foundations, f)
			f.appendCmp2 = dyad.compare_UpSuit
			f.setLabel("A")
		}
		self.tableaux = []*Pile{}
		for x := 0; x < 10; x++ {
			t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
			self.tableaux = append(self.tableaux, t)
			t.appendCmp2 = dyad.compare_DownAltColor
			t.moveCmp2 = dyad.compare_DownAltColor
		}
	case "Easthaven":
		self.waste = nil
		self.foundations = []*Pile{}
		for x := 3; x < 7; x++ {
			f := self.baize.NewFoundation(newPileSlot(x, 0))
			self.foundations = append(self.foundations, f)
			f.appendCmp2 = dyad.compare_UpSuit
			f.setLabel("A")
		}
		self.tableaux = []*Pile{}
		for x := 0; x < 7; x++ {
			t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
			self.tableaux = append(self.tableaux, t)
			t.appendCmp2 = dyad.compare_DownAltColor
			t.moveCmp2 = dyad.compare_DownAltColor
			t.setLabel("K")
		}
	}
}

func (self *Westcliff) StartGame() {
	switch self.variant {
	case "Classic":
		if c := self.stock.extract(0, 1, cardid.CLUB); c != nil {
			self.foundations[0].push(c)
		}
		if c := self.stock.extract(0, 1, cardid.DIAMOND); c != nil {
			self.foundations[1].push(c)
		}
		if c := self.stock.extract(0, 1, cardid.HEART); c != nil {
			self.foundations[2].push(c)
		}
		if c := self.stock.extract(0, 1, cardid.SPADE); c != nil {
			self.foundations[3].push(c)
		}
		fallthrough
	case "American", "Easthaven":
		for _, pile := range self.tableaux {
			for i := 0; i < 2; i++ {
				card := moveCard(self.stock, pile)
				card.flipDown()
			}
		}
		for _, pile := range self.tableaux {
			moveCard(self.stock, pile)
		}
		if self.waste != nil {
			moveCard(self.stock, self.waste)
		}
	}
	self.baize.setRecycles(0)
}

func (self *Westcliff) AfterMove() {
	if self.waste != nil {
		if self.waste.Len() == 0 && self.stock.Len() != 0 {
			moveCard(self.stock, self.waste)
		}
	}
}

func (*Westcliff) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *Westcliff) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (self *Westcliff) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (self *Westcliff) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		switch self.variant {
		case "Classic", "American":
			moveCard(self.stock, self.waste)
		case "Easthaven":
			for _, pile := range self.tableaux {
				moveCard(self.stock, pile)
			}
		}
	} else {
		pile.vtable.tailTapped(tail)
	}
}

// func (*Westcliff) PileTapped(pile *Pile) {}
