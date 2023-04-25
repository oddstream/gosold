package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
)

type Reserve struct {
	pile *Pile
}

func (b *Baize) NewReserve(slot PileSlot, fanType FanType) *Pile {
	pile := b.newPile("Reserve", slot, fanType, MOVE_ONE)
	pile.vtable = &Reserve{pile: pile}
	return pile
}

// canSubtypeAppendTail does some obvious checks on the tail before passing it to the script
func (*Reserve) canSubtypeAppendTail(tail []*Card) (bool, error) {
	return false, errors.New("Cannot add a card to a Reserve")
}

func (self *Reserve) tailTapped(tail []*Card) {
	self.pile.defaultTailTapped(tail)
}

// Conformant when contains zero or one card(s), same as Waste
func (self *Reserve) conformant() bool {
	return self.pile.Len() < 2
}

// unsortedPairs - cards in a reserve pile are always considered to be unsorted
func (self *Reserve) unsortedPairs() int {
	if self.pile.Empty() {
		return 0
	}
	return self.pile.Len() - 1
}

func (self *Reserve) movableTails() [][]*Card {
	return self.pile.singleCardMovableTails()
}
