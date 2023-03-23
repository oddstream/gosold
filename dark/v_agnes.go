package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"

	"oddstream.games/gosold/util"
)

type Agnes struct {
	scriptBase
}

func (self *Agnes) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{0, 0})
	self.waste = nil

	self.foundations = nil
	for x := 3; x < 7; x++ {
		f := self.baize.NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
	}

	self.reserves = nil
	for x := 0; x < 7; x++ {
		r := self.baize.NewReserve(image.Point{x, 1}, FAN_NONE)
		self.reserves = append(self.reserves, r)
	}

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Agnes) StartGame() {

	for _, pile := range self.reserves {
		moveCard(self.stock, pile)
	}

	var dealDown int = 0
	for _, pile := range self.tableaux {
		for i := 0; i < dealDown; i++ {
			card := moveCard(self.stock, pile)
			card.flipDown()
		}
		dealDown++
		moveCard(self.stock, pile)
	}

	c := moveCard(self.stock, self.foundations[0])
	ord := c.Ordinal()
	for _, pile := range self.foundations {
		pile.setLabel(util.OrdinalToShortString(ord))
	}
	ord -= 1
	if ord == 0 {
		ord = 13
	}
	for _, pile := range self.tableaux {
		pile.setLabel(util.OrdinalToShortString(ord))
	}
}

func (self *Agnes) TailMoveError(tail []*Card) (bool, error) {
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

func (self *Agnes) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return self.TwoCards(dst, dst.peek(), tail[0])
}

func (*Agnes) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	switch pile.vtable.(type) {
	case *Foundation:
		return cardPair{c1, c2}.compare_UpSuitWrap()
	case *Tableau:
		return cardPair{c1, c2}.compare_DownAltColorWrap()
	}
	return true, nil
}

func (self *Agnes) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		for _, pile := range self.reserves {
			moveCard(self.stock, pile)
		}
	} else {
		pile.vtable.TailTapped(tail)
	}
}

// func (*Agnes) PileTapped(*Pile) {}
