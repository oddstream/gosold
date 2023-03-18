package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"

	"oddstream.games/gosold/cardid"
)

type Westcliff struct {
	scriptBase
	variant string
}

func (self *Westcliff) BuildPiles() {
	self.stock = self.baize.NewStock(image.Point{0, 0})
	switch self.variant {
	case "Classic":
		self.waste = self.baize.NewWaste(image.Point{1, 0}, FAN_RIGHT3)
		self.foundations = []*Pile{}
		for x := 3; x < 7; x++ {
			f := self.baize.NewFoundation(image.Point{x, 0})
			self.foundations = append(self.foundations, f)
			f.setLabel("A")
		}
		self.tableaux = []*Pile{}
		for x := 0; x < 7; x++ {
			t := self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
			self.tableaux = append(self.tableaux, t)
		}
	case "American":
		self.waste = self.baize.NewWaste(image.Point{1, 0}, FAN_RIGHT3)
		self.foundations = []*Pile{}
		for x := 6; x < 10; x++ {
			f := self.baize.NewFoundation(image.Point{x, 0})
			self.foundations = append(self.foundations, f)
			f.setLabel("A")
		}
		self.tableaux = []*Pile{}
		for x := 0; x < 10; x++ {
			t := self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
			self.tableaux = append(self.tableaux, t)
		}
	case "Easthaven":
		self.waste = nil
		self.foundations = []*Pile{}
		for x := 3; x < 7; x++ {
			f := self.baize.NewFoundation(image.Point{x, 0})
			self.foundations = append(self.foundations, f)
			f.setLabel("A")
		}
		self.tableaux = []*Pile{}
		for x := 0; x < 7; x++ {
			t := self.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
			self.tableaux = append(self.tableaux, t)
			t.setLabel("K")
		}
	}
}

func (self *Westcliff) StartGame() {
	switch self.variant {
	case "Classic":
		if c := self.stock.extract(0, 1, cardid.CLUB); c != nil {
			self.foundations[0].push(c)
		}
		if c := self.stock.extract(0, 1, cardid.DIAMOND); c != nil {
			self.foundations[1].push(c)
		}
		if c := self.stock.extract(0, 1, cardid.HEART); c != nil {
			self.foundations[2].push(c)
		}
		if c := self.stock.extract(0, 1, cardid.SPADE); c != nil {
			self.foundations[3].push(c)
		}
		fallthrough
	case "American", "Easthaven":
		for _, pile := range self.tableaux {
			for i := 0; i < 2; i++ {
				card := moveCard(self.stock, pile)
				card.flipDown()
			}
		}
		for _, pile := range self.tableaux {
			moveCard(self.stock, pile)
		}
		if self.waste != nil {
			moveCard(self.stock, self.waste)
		}
	}
	self.baize.setRecycles(0)
}

func (self *Westcliff) AfterMove() {
	if self.waste != nil {
		if self.waste.Len() == 0 && self.stock.Len() != 0 {
			moveCard(self.stock, self.waste)
		}
	}
}

func (*Westcliff) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := tailConformant(tail, cardPair.compare_DownAltColor)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (*Westcliff) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return compare_Empty(dst, tail[0])
		} else {
			return cardPair{dst.peek(), tail[0]}.compare_UpSuit()
		}
	case *Tableau:
		if dst.Empty() {
			return compare_Empty(dst, tail[0])
		} else {
			return cardPair{dst.peek(), tail[0]}.compare_DownAltColor()
		}
	}
	return true, nil
}

func (*Westcliff) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.compare_DownAltColor)
}

func (self *Westcliff) TailTapped(tail []*Card, nTarget int) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		switch self.variant {
		case "Classic", "American":
			moveCard(self.stock, self.waste)
		case "Easthaven":
			for _, pile := range self.tableaux {
				moveCard(self.stock, pile)
			}
		}
	} else {
		pile.vtable.TailTapped(tail, nTarget)
	}
}

// func (*Westcliff) PileTapped(pile *Pile) {}
