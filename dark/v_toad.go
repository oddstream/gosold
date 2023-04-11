package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"

	"oddstream.games/gosold/util"
)

type Toad struct {
	scriptBase
}

func (self *Toad) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))
	self.waste = self.baize.NewWaste(newPileSlot(1, 0), FAN_RIGHT3)

	self.reserves = nil
	self.reserves = append(self.reserves, self.baize.NewReserve(newPileSlot(3, 0), FAN_RIGHT))

	self.foundations = nil
	for x := 0; x < 8; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 1))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuitWrap
	}

	self.tableaux = nil
	for x := 0; x < 8; x++ {
		// When moving tableau piles, you must either move the whole pile or only the top card.
		t := self.baize.NewTableau(newPileSlot(x, 2), FAN_DOWN, MOVE_ONE_OR_ALL)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownSuitWrap
	}
}

func (self *Toad) StartGame() {
	for n := 0; n < 20; n++ {
		moveCard(self.stock, self.reserves[0])
		self.reserves[0].peek().flipDown()
	}
	self.reserves[0].peek().flipUp()

	for _, pile := range self.tableaux {
		moveCard(self.stock, pile)
	}
	// One card is dealt onto the first foundation. This rank will be used as a base for the other foundations.
	c := moveCard(self.stock, self.foundations[0])
	for _, pile := range self.foundations {
		pile.setLabel(util.OrdinalToShortString(c.Ordinal()))
	}
	moveCard(self.stock, self.waste)
	self.baize.setRecycles(1)
}

func (self *Toad) AfterMove() {
	// Empty spaces are filled automatically from the reserve.
	for _, p := range self.tableaux {
		if p.Empty() {
			moveCard(self.reserves[0], p)
		}
	}
	if self.waste.Len() == 0 && self.stock.Len() != 0 {
		moveCard(self.stock, self.waste)
	}

}

func (*Toad) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (self *Toad) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	card := tail[0]
	if dst.Empty() {
		switch dst.vtable.(type) {
		case *Foundation:
			return compare_Empty(dst, card)
		case *Tableau:
			// Once the reserve is empty, spaces in the tableau can be filled with a card from the Deck [Stock/Waste], but NOT from another tableau pile.
			// pointless rule, since tableuax move rule is MOVE_ONE_OR_ALL
			if card.owner() != self.waste {
				return false, errors.New("Empty tableaux must be filled with cards from the waste")
			}
			return compare_Empty(dst, card)
		}
	}
	return self.TwoCards(dst, dst.peek(), card)
}

func (*Toad) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (self *Toad) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		c := pile.pop()
		self.waste.push(c)
	} else {
		pile.vtable.tailTapped(tail)
	}
}

func (self *Toad) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.waste, self.stock)
	}
}
