package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
)

type Cell struct {
	pile *Pile
}

func (b *Baize) NewCell(slot PileSlot) *Pile {
	pile := b.newPile("Cell", slot, FAN_NONE, MOVE_ONE)
	pile.maxLen = 1
	pile.vtable = &Cell{pile: pile}
	return pile
}

func (self *Cell) canSubtypeAppendTail(tail []*Card) (bool, error) {
	if !self.pile.Empty() {
		return false, errors.New("A Cell can only contain one card")
	}
	if len(tail) > 1 {
		return false, errors.New("Cannot move more than one card to a Cell")
	}
	if anyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card")
	}
	return self.pile.baize.script.TailAppendError(self.pile, tail)
}

func (self *Cell) tailTapped(tail []*Card) {
	self.pile.defaultTailTapped(tail)
}

func (*Cell) conformant() bool {
	return true
}

func (*Cell) unsortedPairs() int {
	return 0
}

func (self *Cell) movableTails() [][]*Card {
	return self.pile.singleCardMovableTails()
}
