package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"

	"oddstream.games/gosold/cardid"
	"oddstream.games/gosold/util"
)

type Canfield struct {
	scriptBase
	variant        string
	draw, recycles int
}

func (self *Canfield) BuildPiles() {

	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	self.wastes = append(self.wastes, self.baize.NewWaste(newPileSlot(1, 0), FAN_RIGHT3))

	self.reserves = append(self.reserves, self.baize.NewReserve(newPileSlot(0, 1), FAN_DOWN))

	for x := 3; x < 7; x++ {
		f := self.baize.NewFoundation(newPileSlot(x, 0))
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuitWrap
	}

	for x := 3; x < 7; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ONE_OR_ALL)
		self.tableaux = append(self.tableaux, t)
		if self.variant == "storehouse" {
			t.appendCmp2 = dyad.compare_DownSuitWrap
			t.moveCmp2 = dyad.compare_DownSuitWrap
		} else {
			t.appendCmp2 = dyad.compare_DownAltColorWrap
			t.moveCmp2 = dyad.compare_DownAltColorWrap
		}
	}
}

func (self *Canfield) StartGame() {

	if self.variant == "storehouse" {
		if c := self.stock.extract(0, 2, cardid.CLUB); c != nil {
			self.foundations[0].push(c)
		}
		if c := self.stock.extract(0, 2, cardid.DIAMOND); c != nil {
			self.foundations[1].push(c)
		}
		if c := self.stock.extract(0, 2, cardid.HEART); c != nil {
			self.foundations[2].push(c)
		}
		if c := self.stock.extract(0, 2, cardid.SPADE); c != nil {
			self.foundations[3].push(c)
		}
	} else if self.variant == "selective" {
		self.baize.fnNotify(MessageEvent, "Move a Tableaux card to a Foundation")
	} else {
		card := moveCard(self.stock, self.foundations[0])
		for _, pile := range self.foundations {
			pile.setLabel(util.OrdinalToShortString(card.Ordinal()))
		}
	}

	for i := 0; i < 12; i++ {
		moveCard(self.stock, self.reserves[0]).flipDown()
	}
	moveCard(self.stock, self.reserves[0]) // face up

	for _, pile := range self.tableaux {
		moveCard(self.stock, pile)
	}

	self.baize.setRecycles(self.recycles)
}

func (self *Canfield) AfterMove() {
	if self.foundations[0].label == "" {
		var ord int = 0
		for _, f := range self.foundations {
			// find where the first card landed
			if len(f.cards) > 0 {
				ord = f.peek().id.Ordinal()
				break
			}
		}
		if ord == 0 {
			self.baize.fnNotify(MessageEvent, "Move a Tableaux card to a Foundation")
		} else {
			for _, f := range self.foundations {
				f.setLabel(util.OrdinalToShortString(ord))
			}
		}
	}
	// "fill each [tableau] space at once with the top card of the reserve,
	// after the reserve is exhausted, fill spaces from the waste pile,
	// but at this time a space may be kept open for as long as desired"
	for _, pile := range self.tableaux {
		if pile.Empty() {
			moveCard(self.reserves[0], pile)
		}
	}
}

// default TailMoveError

func (self *Canfield) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// The top cards are available for play on foundations, BUT NEVER INTO SPACES
	// One card can be moved at a time, but sequences can also be moved as one unit.
	if dst.Empty() {
		switch dst.vtable.(type) {
		case *Foundation:
			if dst.Label() == "" {
				if self.variant == "selective" {
					if tail[0].owner().category != "Tableau" {
						return false, errors.New("The first Foundation card must come from a Tableau")
					}
				}
			}
		case *Tableau:
			// Spaces that occur on the tableau are filled only from reserve or waste
			if tail[0].owner().category == "Tableau" {
				return false, errors.New("An empty Tableau must be filled from the Reserve or Waste")
			}
		}
		return compare_Empty(dst, tail[0])
	}
	return dst.appendCmp2(dyad{dst.peek(), tail[0]})
}

func (self *Canfield) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		for i := 0; i < self.draw; i++ {
			moveCard(self.stock, self.Waste())
		}
	} else {
		pile.vtable.tailTapped(tail)
	}
}

func (self *Canfield) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.Waste(), self.stock)
	}
}
