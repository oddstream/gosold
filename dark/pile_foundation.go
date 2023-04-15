package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
)

type Foundation struct {
	pile *Pile
}

func (b *Baize) NewFoundation(slot PileSlot) *Pile {
	pile := b.newPile("Foundation", slot, FAN_NONE, MOVE_NONE)
	pile.vtable = &Foundation{pile: pile}
	return pile
}

// canSubtypeAppendTail does some obvious checks on the tail before passing it to the script
func (self *Foundation) canSubtypeAppendTail(tail []*Card) (bool, error) {
	if len(tail) > 1 {
		return false, errors.New("Cannot move more than one card to a Foundation")
	}
	if self.pile.Len() == 13 {
		return false, errors.New("That Foundation already contains 13 cards")
	}
	if anyCardsProne(tail) {
		return false, errors.New("Cannot add a face down card to a Foundation")
	}
	return self.pile.baize.script.TailAppendError(self.pile, tail)
}

func (*Foundation) tailTapped([]*Card) {}

func (*Foundation) conformant() bool {
	return true
}

func (*Foundation) unsortedPairs() int {
	// you can only put a sorted sequence into a Foundation, so this will always be zero
	return 0
}

func (*Foundation) movableTails() [][]*Card {
	return nil
}
