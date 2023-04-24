package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"

	"oddstream.games/gosold/util"
)

type Duchess struct {
	scriptBase
}

func (self *Duchess) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(1, 1))

	for i := 0; i < 4; i++ {
		r := self.baize.NewReserve(newPileSlot(i*2, 0), FAN_RIGHT)
		self.reserves = append(self.reserves, r)
	}

	self.wastes = append(self.wastes, self.baize.NewWaste(newPileSlot(1, 2), FAN_DOWN3))

	for x := 3; x < 7; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 1))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuitWrap
	}

	for x := 3; x < 7; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 2), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_DownAltColorWrap
		t.moveCmp2 = dyad.compare_DownAltColorWrap
	}
}

func (self *Duchess) StartGame() {
	for _, pile := range self.foundations {
		pile.setLabel("")
	}
	for _, pile := range self.reserves {
		moveCard(self.stock, pile)
		moveCard(self.stock, pile)
		moveCard(self.stock, pile)
	}
	for _, pile := range self.tableaux {
		moveCard(self.stock, pile)
	}
	self.baize.setRecycles(1)
	self.baize.fnNotify(MessageEvent, "Move a Reserve card to a Foundation")
}

func (self *Duchess) AfterMove() {
	if self.foundations[0].label == "" {
		// To start the game, the player will choose among the top cards of the reserve fans which will start the first foundation pile.
		// Once he/she makes that decision and picks a card, the three other cards with the same rank,
		// whenever they become available, will start the other three foundations.
		var ord int = 0
		for _, f := range self.foundations {
			// find where the first card landed
			if len(f.cards) > 0 {
				ord = f.peek().id.Ordinal()
				break
			}
		}
		if ord == 0 {
			self.baize.fnNotify(MessageEvent, "Move a Reserve card to a Foundation")
		} else {
			for _, f := range self.foundations {
				f.setLabel(util.OrdinalToShortString(ord))
			}
		}
	}
}

// default TailMoveError

func (self *Duchess) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		switch dst.vtable.(type) {
		case *Foundation:
			if dst.Label() == "" {
				if tail[0].owner().category != "Reserve" {
					return false, errors.New("The first Foundation card must come from a Reserve")
				}
			}
		case *Tableau:
			var rescards int = 0
			for _, p := range self.reserves {
				rescards += p.Len()
			}
			if rescards > 0 {
				// Spaces that occur on the tableau are filled with any top card in the reserve
				if tail[0].owner().category != "Reserve" {
					return false, errors.New("An empty Tableau must be filled from a Reserve")
				}
			}
		}
		return compare_Empty(dst, tail[0])
	}
	return dst.appendCmp2(dyad{dst.peek(), tail[0]})
}

// default TailTapped

func (self *Duchess) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.Waste(), self.stock)
	}
}
