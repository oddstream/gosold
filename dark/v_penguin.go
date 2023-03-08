package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"

	"oddstream.games/gosold/util"
)

type Penguin struct {
	scriptBase
}

func (pen *Penguin) BuildPiles() {

	// hidden (off-screen) stock
	pen.stock = pen.baize.NewStock(image.Point{-5, -5}, FAN_NONE, 1, 4, nil, 0)
	pen.waste = nil

	// the flipper, seven cells
	pen.cells = nil
	for x := 0; x < 7; x++ {
		pile := pen.baize.NewCell(image.Point{x, 0})
		pen.cells = append(pen.cells, pile)
	}

	pen.foundations = nil
	for y := 0; y < 4; y++ {
		pile := pen.baize.NewFoundation(image.Point{8, y})
		pen.foundations = append(pen.foundations, pile)
	}

	pen.tableaux = nil
	for x := 0; x < 7; x++ {
		t := pen.baize.NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		pen.tableaux = append(pen.tableaux, t)
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
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := tailConformant(tail, cardPair.Compare_DownSuitWrap)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (*Penguin) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return compare_Empty(dst, tail[0])
		} else {
			return cardPair{dst.peek(), tail[0]}.compare_UpSuitWrap()
		}
	case *Tableau:
		if dst.Empty() {
			return compare_Empty(dst, tail[0])
		} else {
			return cardPair{dst.peek(), tail[0]}.Compare_DownSuitWrap()
		}
	}
	return true, nil
}

func (*Penguin) UnsortedPairs(pile *Pile) int {
	return unsortedPairs(pile, cardPair.Compare_DownSuitWrap)
}

func (pen *Penguin) TailTapped(tail []*Card) {
	tail[0].owner().vtable.TailTapped(tail)
}

// func (pen *Penguin) PileTapped(pile *Pile) {}
