package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"

	"oddstream.games/gosold/cardid"
)

type FortyThieves struct {
	scriptBase
	founds         []int
	tabs           []int
	proneRows      []int
	cardsPerTab    int
	recycles       int
	dealAces       bool
	moveType       MoveType
	tabCompareFunc cardPairCompareFunc
}

func (self *FortyThieves) BuildPiles() {

	if self.moveType == MOVE_NONE /* 0 */ {
		self.moveType = MOVE_ONE_PLUS
	}
	if self.cardColors == 0 {
		self.cardColors = 2
	}
	if self.tabCompareFunc == nil {
		self.tabCompareFunc = cardPair.compare_DownSuit
	}

	self.stock = self.baize.NewStock(image.Point{0, 0})
	self.waste = self.baize.NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	self.foundations = nil
	for _, x := range self.founds {
		f := self.baize.NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.setLabel("A")
	}

	self.tableaux = nil
	for _, x := range self.tabs {
		t := self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, self.moveType)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *FortyThieves) StartGame() {
	if self.dealAces {
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
		if c := self.stock.extract(1, 1, cardid.CLUB); c != nil {
			self.foundations[4].push(c)
		}
		if c := self.stock.extract(1, 1, cardid.DIAMOND); c != nil {
			self.foundations[5].push(c)
		}
		if c := self.stock.extract(1, 1, cardid.HEART); c != nil {
			self.foundations[6].push(c)
		}
		if c := self.stock.extract(1, 1, cardid.SPADE); c != nil {
			self.foundations[7].push(c)
		}
	}
	for _, pile := range self.tableaux {
		for i := 0; i < self.cardsPerTab; i++ {
			moveCard(self.stock, pile)
		}
	}
	for _, row := range self.proneRows {
		for _, pile := range self.tableaux {
			pile.cards[row].flipDown()
		}
	}
	self.baize.setRecycles(self.recycles)
	moveCard(self.stock, self.waste)
}

func (self *FortyThieves) AfterMove() {
	if self.waste.Empty() && !self.stock.Empty() {
		moveCard(self.stock, self.waste)
	}
}

func (self *FortyThieves) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := tailConformant(tail, self.tabCompareFunc)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (self *FortyThieves) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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
			return self.tabCompareFunc(cardPair{dst.peek(), tail[0]})
		}
	}
	return true, nil
}

func (self *FortyThieves) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, self.tabCompareFunc)
}

func (self *FortyThieves) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		moveCard(self.stock, self.waste)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *FortyThieves) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.waste, self.stock)
	}
}
