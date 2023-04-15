package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
)

type Waste struct {
	pile *Pile
}

func (b *Baize) NewWaste(slot PileSlot, fanType FanType) *Pile {
	pile := b.newPile("Waste", slot, fanType, MOVE_ONE)
	pile.vtable = &Waste{pile: pile}
	pile.appendFrom = "Stock"
	return pile
}

func (self *Waste) canSubtypeAppendTail(tail []*Card) (bool, error) {
	if len(tail) > 1 {
		return false, errors.New("Can only move a single card to Waste")
	}
	return self.pile.baize.script.TailAppendError(self.pile, tail)
}

func (self *Waste) tailTapped(tail []*Card) {
	self.pile.defaultTailTapped(tail)
}

// Conformant when contains zero or one card(s), same as Reserve
func (self *Waste) conformant() bool {
	return self.pile.Len() < 2
}

// unsortedPairs - cards in a waste pile are always considered to be unsorted
func (self *Waste) unsortedPairs() int {
	if self.pile.Empty() {
		return 0
	}
	return self.pile.Len() - 1
}

func (self *Waste) movableTails() [][]*Card {
	return self.pile.singleCardMovableTails()
}
