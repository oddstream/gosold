package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"fmt"
	"image"

	"oddstream.games/gosold/cardid"
	"oddstream.games/gosold/util"
)

type Canfield struct {
	scriptBase
	variant        string
	draw, recycles int
	tabCompareFunc cardPairCompareFunc
}

func (self *Canfield) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{0, 0})

	self.waste = self.baize.NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	self.reserves = nil
	self.reserves = append(self.reserves, self.baize.NewReserve(image.Point{0, 1}, FAN_DOWN))

	self.foundations = nil
	for x := 3; x < 7; x++ {
		self.foundations = append(self.foundations, self.baize.NewFoundation(image.Point{x, 0}))
	}

	self.tableaux = nil
	for x := 3; x < 7; x++ {
		self.tableaux = append(self.tableaux, self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_OR_ALL))
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
	} else {
		card := moveCard(self.stock, self.foundations[0])
		for _, pile := range self.foundations {
			pile.setLabel(util.OrdinalToShortString(card.Ordinal()))
		}
	}

	for i := 0; i < 12; i++ {
		moveCard(self.stock, self.reserves[0]).flipDown()
	}
	moveCard(self.stock, self.reserves[0])

	for _, pile := range self.tableaux {
		moveCard(self.stock, pile)
	}

	self.baize.setRecycles(self.recycles)
}

func (self *Canfield) AfterMove() {
	// "fill each [tableau] space at once with the top card of the reserve,
	// after the reserve is exhausted, fill spaces from the waste pile,
	// but at this time a space may be kept open for as long as desired"
	for _, pile := range self.tableaux {
		if pile.Empty() {
			moveCard(self.reserves[0], pile)
		}
	}
}

func (self *Canfield) TailMoveError(tail []*Card) (bool, error) {
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

func (self *Canfield) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// The top cards are available for play on foundations, BUT NEVER INTO SPACES
	// One card can be moved at a time, but sequences can also be moved as one unit.
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			c := tail[0]
			ord := util.OrdinalToShortString(c.Ordinal())
			if dst.Label() == "" {
				if c.owner().category != "Reserve" {
					return false, errors.New("The first Foundation card must come from a Reserve")
				}
				for _, pile := range self.foundations {
					pile.setLabel(ord)
				}
			}
			if ord != dst.Label() {
				return false, fmt.Errorf("Foundations can only accept an %s, not a %s", dst.Label(), ord)
			}
		} else {
			return cardPair{dst.peek(), tail[0]}.compare_UpSuitWrap()
		}
	case *Tableau:
		if dst.Empty() {
			// Spaces that occur on the tableau are filled only from reserve or waste
			if tail[0].owner().category == "Tableau" {
				return false, errors.New("An empty Tableau must be filled from the Reserve or Waste")
			}
			return true, nil
		} else {
			return self.tabCompareFunc(cardPair{dst.peek(), tail[0]})
		}
	}
	return true, nil
}

func (self *Canfield) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, self.tabCompareFunc)
}

func (self *Canfield) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		for i := 0; i < self.draw; i++ {
			moveCard(self.stock, self.waste)
		}
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *Canfield) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.waste, self.stock)
	}
}
