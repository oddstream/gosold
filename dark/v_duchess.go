package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"

	"oddstream.games/gosold/util"
)

type Duchess struct {
	scriptBase
}

func (self *Duchess) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{1, 1}, FAN_NONE, 1, 4, nil, 0)

	self.reserves = []*Pile{}
	for i := 0; i < 4; i++ {
		self.reserves = append(self.reserves, self.baize.NewReserve(image.Point{i * 2, 0}, FAN_RIGHT))
	}

	self.waste = self.baize.NewWaste(image.Point{1, 2}, FAN_DOWN3)

	self.foundations = []*Pile{}
	for x := 3; x < 7; x++ {
		self.foundations = append(self.foundations, self.baize.NewFoundation(image.Point{x, 1}))
	}

	self.tableaux = []*Pile{}
	for x := 3; x < 7; x++ {
		self.tableaux = append(self.tableaux, self.baize.NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ANY))
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

		} else {
			for _, f := range self.foundations {
				f.setLabel(util.OrdinalToShortString(ord))
			}
		}
	}
}

func (*Duchess) TailMoveError(tail []*Card) (bool, error) {
	// One card can be moved at a time, but sequences can also be moved as one unit.
	var pile *Pile = tail[0].owner()
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := tailConformant(tail, cardPair.compare_DownAltColorWrap)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (self *Duchess) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	card := tail[0]
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			if dst.Label() == "" {
				if card.owner().category != "Reserve" {
					return false, errors.New("The first Foundation card must come from a Reserve")
				}
			}
			return compare_Empty(dst, card)
		} else {
			return cardPair{dst.peek(), card}.compare_UpSuitWrap()
		}
	case *Tableau:
		if dst.Empty() {
			var rescards int = 0
			for _, p := range self.reserves {
				rescards += p.Len()
			}
			if rescards > 0 {
				// Spaces that occur on the tableau are filled with any top card in the reserve
				if card.owner().category != "Reserve" {
					return false, errors.New("An empty Tableau must be filled from a Reserve")
				}
			}
			return true, nil
		} else {
			return cardPair{dst.peek(), card}.compare_DownAltColorWrap()
		}
	}
	return true, nil
}

func (*Duchess) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.compare_DownAltColorWrap)
}

func (self *Duchess) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		moveCard(self.stock, self.waste)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *Duchess) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.waste, self.stock)
	}
}
