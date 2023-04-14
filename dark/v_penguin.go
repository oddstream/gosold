package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"oddstream.games/gosold/util"
)

type Penguin struct {
	scriptBase
}

func (pen *Penguin) BuildPiles() {

	// hidden (off-screen) stock
	pen.stock = pen.baize.NewStock(newHiddenPileSlot())

	// the flipper, seven cells
	pen.cells = nil
	for x := 0; x < 7; x++ {
		c := pen.baize.NewCell(newPileSlot(x, 0))
		pen.cells = append(pen.cells, c)
	}

	pen.foundations = nil
	for y := 0; y < 4; y++ {
		f := pen.baize.NewFoundation(newPileSlot(8, y))
		pen.foundations = append(pen.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuitWrap
	}

	pen.tableaux = nil
	for x := 0; x < 7; x++ {
		t := pen.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		pen.tableaux = append(pen.tableaux, t)
		t.appendCmp2 = dyad.compare_DownSuitWrap
		t.moveCmp2 = dyad.compare_DownSuitWrap
	}
}

func (pen *Penguin) StartGame() {
	// Shuffle a 52-card pack and deal the first card face up to the top left of the board.
	// This card is called the Beak.

	beak := moveCard(pen.stock, pen.tableaux[0])
	for _, pile := range pen.foundations {
		pile.setLabel(util.OrdinalToShortString(beak.Ordinal()))
	}

	var fnext int = 0 // the next foundation we will move a card to

	// 49-card layout consisting of seven rows and seven columns
	for _, pile := range pen.tableaux {
		for pile.Len() < 7 {
			// As and when the other three cards of the same rank turn up in the deal,
			// take them out and set them apart as foundations.
			card := pen.stock.peek()
			if card.Ordinal() == beak.Ordinal() {
				moveCard(pen.stock, pen.foundations[fnext])
				fnext += 1
			} else {
				moveCard(pen.stock, pile)
			}
		}
	}
	for pen.stock.Len() > 0 {
		// we have 7x7 cards in tableaux, remaining cards must be ordinal == beak
		moveCard(pen.stock, pen.foundations[fnext])
		fnext += 1
	}

	// When you empty a column, you may fill the space it leaves with a card one rank lower than the rank of the beak,
	// together with any other cards attached to it in descending suit-sequence.
	// For example, since the beak is a Ten, you can start a new column only with a Nine,
	// or a suit-sequence headed by a Nine.

	ord := beak.Ordinal() - 1
	if ord == 0 {
		ord = 13
	}
	for _, pile := range pen.tableaux {
		pile.setLabel(util.OrdinalToShortString(ord))
	}
}

func (*Penguin) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].owner()
	return tailConformant(tail, pile.moveCmp2)
}

func (pen *Penguin) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	if dst.Empty() {
		return compare_Empty(dst, tail[0])
	}
	return pen.TwoCards(dst, dst.peek(), tail[0])
}

func (*Penguin) TwoCards(pile *Pile, c1, c2 *Card) (bool, error) {
	return pile.appendCmp2(dyad{c1, c2})
}

func (pen *Penguin) TailTapped(tail []*Card) {
	tail[0].owner().vtable.tailTapped(tail)
}

// func (pen *Penguin) PileTapped(pile *Pile) {}
