package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
)

type Discard struct {
	pile *Pile
}

func (b *Baize) NewDiscard(slot PileSlot, fanType FanType) *Pile {
	pile := b.newPile("Discard", slot, FAN_NONE, MOVE_NONE)
	pile.vtable = &Discard{pile: pile}
	pile.appendFrom = "Tableau"
	pile.maxLen = 13 // TODO stripped decks
	return pile
}

// canSubtypeAppendTail does some obvious checks on the tail before passing it to the script
func (self *Discard) canSubtypeAppendTail(tail []*Card) (bool, error) {
	if !self.pile.Empty() {
		return false, errors.New("Can only move cards to an empty Discard")
	}
	if anyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card to a Discard")
	}
	// TODO this should be number of cards in a suit
	// which is usually 13
	// but may be fewer in a stripped deck
	if len(tail) != 13 {
		return false, errors.New("Can only move a full set of cards to a Discard")
	}
	if ok, err := tailConformant(tail, dyad.compare_DownSuit); !ok {
		return false, err
	}
	return self.pile.baize.script.TailAppendError(self.pile, tail)
}

func (*Discard) tailTapped([]*Card) {
	// do nothing
}

func (*Discard) conformant() bool {
	// no Baize that contains any discard piles should be Conformant,
	// because there is no use showing the collect all FAB
	// because that would do nothing
	// because cards are not collected to discard piles
	return false
}

func (*Discard) unsortedPairs() int {
	// you can only put a sequence into a Discard, so this will always be zero
	return 0
}

func (*Discard) movableTails() [][]*Card {
	return nil
}
