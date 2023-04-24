package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"oddstream.games/gosold/cardid"
)

type Alhambra struct {
	scriptBase
}

func (self *Alhambra) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 3))

	// waste pile implemented as a tableau because cards may be built on it
	t := self.baize.NewTableau(newPileSlot(1, 3), FAN_RIGHT3, MOVE_ONE)
	self.tableaux = append(self.tableaux, t)
	t.appendCmp2 = dyad.compare_UpOrDownSuitWrap

	for x := 0; x < 4; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}
	for x := 4; x < 8; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_DownSuit
		f.setLabel("K")
	}

	for x := 0; x < 8; x++ {
		r := self.baize.NewReserve(newPileSlot(x, 1), FAN_DOWN)
		self.reserves = append(self.reserves, r)
	}
}

func (self *Alhambra) StartGame() {

	self.foundations[0].push(self.stock.extract(0, 1, cardid.CLUB))
	self.foundations[1].push(self.stock.extract(0, 1, cardid.DIAMOND))
	self.foundations[2].push(self.stock.extract(0, 1, cardid.HEART))
	self.foundations[3].push(self.stock.extract(0, 1, cardid.SPADE))
	self.foundations[4].push(self.stock.extract(0, 13, cardid.CLUB))
	self.foundations[5].push(self.stock.extract(0, 13, cardid.DIAMOND))
	self.foundations[6].push(self.stock.extract(0, 13, cardid.HEART))
	self.foundations[7].push(self.stock.extract(0, 13, cardid.SPADE))

	for _, r := range self.reserves {
		for i := 0; i < 4; i++ {
			moveCard(self.stock, r)
		}
	}

	self.baize.setRecycles(2)
}

func (*Alhambra) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

// default TailAppendError

// TailTapped override default to move card to tableaux[0]
func (self *Alhambra) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		moveCard(self.stock, self.tableaux[0])
	} else {
		pile.vtable.tailTapped(tail)
	}
}

func (self *Alhambra) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.tableaux[0], self.stock)
	}
}
