package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"log"
)

type Spider struct {
	scriptBase
}

func (self *Spider) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	self.discards = nil
	for x := 2; x < 10; x++ {
		d := self.baize.NewDiscard(newPileSlot(x, 0), FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = nil
	for x := 0; x < 10; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_Down
		t.moveCmp2 = dyad.compare_DownSuit
	}
}

func (self *Spider) StartGame() {
	// The Tableau consists of 10 stacks with 6 cards in the first 4 stacks, with the 6th card face up,
	// and 5 cards in the remaining 6 stacks, with the 5th card face up.

	for i := 0; i < 4; i++ {
		pile := self.tableaux[i]
		for j := 0; j < 6; j++ {
			moveCard(self.stock, pile).flipDown()
		}
	}
	for i := 4; i < 10; i++ {
		pile := self.tableaux[i]
		for j := 0; j < 5; j++ {
			moveCard(self.stock, pile).flipDown()
		}
	}
	for _, pile := range self.tableaux {
		c := pile.peek()
		if c == nil {
			log.Panic("empty tableau")
		}
		c.flipUp()
	}
	self.baize.setRecycles(0)
}

func (*Spider) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (self *Spider) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*Spider) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (self *Spider) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	switch pile.vtable.(type) {
	case *Stock:
		var tabCards, emptyTabs int
		for _, tab := range self.tableaux {
			if tab.Len() == 0 {
				emptyTabs++
			} else {
				tabCards += tab.Len()
			}
		}
		if emptyTabs > 0 && tabCards >= len(self.tableaux) {
			self.baize.fnNotify(MessageEvent, "All empty tableaux must be filled before dealing a new row")
		} else {
			for _, tab := range self.tableaux {
				moveCard(self.stock, tab)
			}
		}
	default:
		tail[0].owner().vtable.tailTapped(tail)
	}
}

// func (*Spider) PileTapped(*Pile) {}

func (self *Spider) Complete() bool {
	return self.SpiderComplete()
}
