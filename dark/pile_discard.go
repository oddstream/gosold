package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
)

type Discard struct {
	pile *Pile
}

func (b *Baize) NewDiscard(slot image.Point, fanType FanType) *Pile {
	pile := b.newPile("Discard", slot, FAN_NONE, MOVE_NONE)
	pile.vtable = &Discard{pile: pile}
	return pile
}

func (self *Discard) CanAcceptTail(tail []*Card) (bool, error) {
	if !self.pile.Empty() {
		return false, errors.New("Can only move cards to an empty Discard")
	}
	if anyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card to a Discard")
	}
	if len(tail) != self.pile.baize.cardCount/len(self.pile.baize.script.Discards()) {
		return false, errors.New("Can only move a full set of cards to a Discard")
	}
	if ok, err := tailConformant(tail, cardPair.compare_DownSuit); !ok {
		return false, err
	}
	// Scorpion tails can always be moved, but Mrs Mop/Simple Simon tails
	// must be conformant, so ...
	return self.pile.baize.script.TailMoveError(tail)
}

func (*Discard) TailTapped([]*Card, int) {
	// do nothing
}

func (*Discard) Conformant() bool {
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

func (*Discard) MovableTails2() [][]*Card {
	return nil
}
