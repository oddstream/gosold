package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
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
	tabCompareFunc dyadCmpFunc
}

func (self *FortyThieves) BuildPiles() {

	if self.moveType == MOVE_NONE /* 0 */ {
		self.moveType = MOVE_ONE_PLUS
	}
	if self.cardColors == 0 {
		self.cardColors = 2
	}
	if self.tabCompareFunc == nil {
		self.tabCompareFunc = dyad.compare_DownSuit
	}

	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	self.wastes = append(self.wastes, self.baize.NewWaste(newPileSlot(1, 0), FAN_RIGHT3))

	for _, x := range self.founds {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmpFunc = dyad.compare_UpSuit
		f.setLabel("A")
	}

	for _, x := range self.tabs {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, self.moveType)
		self.tableaux = append(self.tableaux, t)
		t.appendCmpFunc = self.tabCompareFunc
		t.moveCmpFunc = self.tabCompareFunc
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
	moveCard(self.stock, self.Waste())
}

func (self *FortyThieves) AfterMove() {
	if self.Waste().Empty() && !self.stock.Empty() {
		moveCard(self.stock, self.Waste())
	}
}

// default TailMoveError

// default TailAppendError

// default TailTapped

func (self *FortyThieves) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.Waste(), self.stock)
	}
}
