package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"
	"log"

	"oddstream.games/gosold/util"
)

type CanThieves struct {
	scriptBase
}

func (self *CanThieves) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{0, 0})
	self.waste = self.baize.NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	if self.reserves != nil {
		log.Println("*** reserves is not nil ***")
	}
	self.reserves = nil
	self.reserves = append(self.reserves, self.baize.NewReserve(image.Point{0, 1}, FAN_DOWN))

	self.foundations = nil
	for x := 3; x < 11; x++ {
		self.foundations = append(self.foundations, self.baize.NewFoundation(image.Point{x, 0}))
	}

	self.tableaux = nil
	for x := 2; x < 6; x++ {
		self.tableaux = append(self.tableaux, self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY))
	}
	for x := 7; x < 12; x++ {
		self.tableaux = append(self.tableaux, self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY))
	}
}

func (self *CanThieves) StartGame() {
	for _, pile := range self.foundations {
		pile.setLabel("")
	}

	// "At the start of the game 13 cards are dealt here."
	for i := 0; i < 12; i++ {
		moveCard(self.stock, self.reserves[0]).flipDown()
	}
	moveCard(self.stock, self.reserves[0])

	// "At the start of the game 1 card is dealt each to the left 4 piles,
	// and 8 cards are dealt to each of the remaining 5 piles."
	for i := 0; i < 4; i++ {
		moveCard(self.stock, self.tableaux[i])
	}
	for i := 4; i < 9; i++ {
		for j := 0; j < 8; j++ {
			moveCard(self.stock, self.tableaux[i])
		}
	}

	self.baize.setRecycles(2)
}

func (self *CanThieves) AfterMove() {
	if self.foundations[0].label == "" {
		// The first card played to a foundation will determine the starting ordinal for all the foundations
		var ord int = 0
		for _, f := range self.foundations {
			// find where the first card landed
			if len(f.cards) > 0 {
				ord = f.peek().id.Ordinal()
				break
			}
		}
		if ord != 0 {
			for _, f := range self.foundations {
				f.setLabel(util.OrdinalToShortString(ord))
			}
		}
	}
}

func (self *CanThieves) inFirstFour(tab *Pile) bool {
	for i := 0; i < 4; i++ {
		if tab == self.tableaux[i] {
			return true
		}
	}
	return false
}

func (self *CanThieves) TailMoveError(tail []*Card) (bool, error) {
	// One card can be moved at a time, but sequences can also be moved as one unit.
	var pile *Pile = tail[0].owner()
	switch pile.vtable.(type) {
	case *Tableau:
		if self.inFirstFour(pile) {
			ok, err := tailConformant(tail, cardPair.compare_DownAltColorWrap)
			if !ok {
				return ok, err
			}
		} else {
			ok, err := tailConformant(tail, cardPair.Compare_DownSuitWrap)
			if !ok {
				return ok, err
			}
		}
	}
	return true, nil
}

func (self *CanThieves) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	card := tail[0]
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return compare_Empty(dst, card)
		} else {
			return cardPair{dst.peek(), card}.compare_UpSuitWrap()
		}
	case *Tableau:
		if dst.Empty() {
			if self.inFirstFour(dst) {
				return tailConformant(tail, cardPair.compare_DownAltColorWrap)
			} else {
				return tailConformant(tail, cardPair.Compare_DownSuitWrap)
			}
		} else {
			if self.inFirstFour(dst) {
				ok, err := tailConformant(tail, cardPair.compare_DownAltColorWrap)
				if !ok {
					return ok, err
				}
				return cardPair{dst.peek(), card}.compare_DownAltColorWrap()
			} else {
				ok, err := tailConformant(tail, cardPair.Compare_DownSuitWrap)
				if !ok {
					return ok, err
				}
				return cardPair{dst.peek(), card}.Compare_DownSuitWrap()
			}
		}
	}
	return true, nil
}

func (self *CanThieves) UnsortedPairs(pile *Pile) int {
	switch pile.vtable.(type) {
	case *Tableau:
		if self.inFirstFour(pile) {
			return unsortedPairs(pile, cardPair.compare_DownAltColorWrap)
		} else {
			return unsortedPairs(pile, cardPair.Compare_DownSuitWrap)
		}
	default:
		log.Println("*** eh?", pile.category)
	}
	return 0
}

func (self *CanThieves) TailTapped(tail []*Card, nTarget int) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		moveCard(self.stock, self.waste)
	} else {
		pile.vtable.TailTapped(tail, nTarget)
	}
}

func (self *CanThieves) PileTapped(pile *Pile) {
	if pile == self.stock {
		recycleWasteToStock(self.waste, self.stock)
	}
}
