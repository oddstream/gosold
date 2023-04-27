package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"oddstream.games/gosold/util"
)

type EagleWing struct {
	scriptBase
	trunk *Pile
}

/*
	As described there is no building allowed in the tableau,
	and this makes wins extremely rare.
	In her book 100 Games of Solitaire,
	Helen L. Coops allows building down by suit in the tableau.
	Many software implementations of Eagle Wing follow these rules,
	with a maximum of three cards per space.
	This makes the chances of winning as much as 80%.
*/

func (self *EagleWing) BuildPiles() {

	self.stock = self.baize.NewStock(PileSlot{X: 3.5, Y: 3, Deg: 0})

	self.wastes = append(self.wastes, self.baize.NewWaste(PileSlot{X: 4.5, Y: 3, Deg: 0}, FAN_RIGHT3))

	self.reserves = append(self.reserves, self.baize.NewReserve(newPileSlot(4, 1), FAN_NONE))
	self.trunk = self.reserves[0]

	for x := 4; x < 8; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuitWrap
	}

	for x := 0; x < 9; x++ {
		if x == 4 {
			continue
		}
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownSuitWrap
		t.maxLen = 3
	}
}

func (self *EagleWing) StartGame() {

	for i := 0; i < 13; i++ {
		card := moveCard(self.stock, self.trunk)
		card.flipDown()
	}

	for _, c := range self.tableaux {
		moveCard(self.stock, c)
	}

	card := moveCard(self.stock, self.foundations[0])
	for _, pile := range self.foundations {
		pile.setLabel(util.OrdinalToShortString(card.Ordinal()))
	}

	self.populateWasteFromStock(1)
	self.baize.setRecycles(2)
}

func (self *EagleWing) AfterMove() {

	for _, t := range self.tableaux {
		if t.Empty() {
			moveCard(self.trunk, t)
			if card := self.trunk.peek(); card != nil {
				card.flipDown()
			}
		}
	}

	if self.trunk.Len() == 1 {
		if card := self.trunk.peek(); card != nil {
			card.flipUp()
		}
	}
}

// default TailMoveError

func (self *EagleWing) TailAppendError(dst *Pile, tail []*Card) (bool, error) {

	// Foundation cards must come from the wings, except if there is only one trunk card remaining
	// switch dst.vtable.(type) {
	// case *Foundation:
	// 	src := tail[0].owner()
	// 	if src == self.Waste() {
	// 		return false, errors.New("Foundation cards must come from the wings")
	// 	}
	// 	// if self.trunk.Len() == 1 {
	// 	// }
	// }

	if dst.Empty() {
		return compare_Empty(dst, tail)
	}
	return dst.appendCmp2(dyad{dst.peek(), tail[0]})
}

// default TailTapped

func (self *EagleWing) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.Waste(), self.stock)
	}
}
