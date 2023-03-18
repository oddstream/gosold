package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"

	"oddstream.games/gosold/cardid"
)

type Bisley struct {
	scriptBase
}

func (self *Bisley) BuildPiles() {

	self.stock = self.baize.NewStock(image.Point{-5, -5})

	self.foundations = nil

	for x := 0; x < 4; x++ {
		f := self.baize.NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.setLabel("K")
	}

	for x := 0; x < 4; x++ {
		f := self.baize.NewFoundation(image.Point{x, 1})
		self.foundations = append(self.foundations, f)
		f.setLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 13; x++ {
		t := self.baize.NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		t.setLabel("X")
	}
}

func (self *Bisley) StartGame() {

	self.foundations[4].push(self.stock.extract(0, 1, cardid.CLUB))
	self.foundations[5].push(self.stock.extract(0, 1, cardid.DIAMOND))
	self.foundations[6].push(self.stock.extract(0, 1, cardid.HEART))
	self.foundations[7].push(self.stock.extract(0, 1, cardid.SPADE))

	// the first 4 tableaux have 3 cards
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			moveCard(self.stock, self.tableaux[i])
		}
	}
	// the next 9 tableaux have 4 cards
	for i := 4; i < 13; i++ {
		for j := 0; j < 4; j++ {
			moveCard(self.stock, self.tableaux[i])
		}
	}

	self.baize.setRecycles(0)
}

func (*Bisley) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (self *Bisley) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return compare_Empty(dst, tail[0])
		} else {
			if dst.Label() == "A" {
				return cardPair{dst.peek(), tail[0]}.compare_UpSuit()
			} else {
				return cardPair{dst.peek(), tail[0]}.compare_DownSuit()
			}
		}
	case *Tableau:
		if dst.Empty() {
			return compare_Empty(dst, tail[0])
		} else {
			// return CardPair{dst.peek(), tail[0]}.Compare_UpOrDownSuit()
			return cardPair{dst.peek(), tail[0]}.chainCall(cardPair.compare_UpOrDown, cardPair.compare_Suit)
		}
	}
	return true, nil
}

func (*Bisley) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.compare_DownColor)
}

func (self *Bisley) TailTapped(tail []*Card, nTarget int) {
	var pile *Pile = tail[0].owner()
	if pile == self.stock && len(tail) == 1 {
		moveCard(self.stock, self.waste)
	} else {
		pile.vtable.TailTapped(tail, nTarget)
	}
}

// func (*Bisley) PileTapped(*Pile) {}
