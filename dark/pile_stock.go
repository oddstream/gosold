package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
)

type Stock struct {
	pile *Pile
}

func (b *Baize) NewStock(slot PileSlot) *Pile {
	pile := b.newPile("Stock", slot, FAN_NONE, MOVE_ONE)
	pile.vtable = &Stock{pile: pile}
	return pile
}

func (*Stock) canSubtypeAppendTail([]*Card) (bool, error) {
	return false, errors.New("Cannot move cards to the Stock")
}

func (*Stock) tailTapped([]*Card) {
	// do nothing, handled by script, which had first dibs
}

func (self *Stock) conformant() bool {
	return self.pile.Empty()
}

// unsortedPairs - cards in a stock pile are always considered to be unsorted
func (self *Stock) unsortedPairs() int {
	if self.pile.Empty() {
		return 0
	}
	return self.pile.Len() - 1
}

func (self *Stock) movableTails() [][]*Card {
	return self.pile.singleCardMovableTails()
}
