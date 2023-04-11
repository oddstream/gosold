package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"

	"oddstream.games/gosold/util"
)

type cardPair struct {
	c1, c2 *Card
}

type cardPairCompareFunc func(cardPair) (bool, error)

func tailConformant(tail []*Card, fn cardPairCompareFunc) (bool, error) {
	for i := 1; i < len(tail); i++ {
		if ok, err := fn(cardPair{tail[i-1], tail[i]}); !ok {
			return false, err
		}
	}
	return true, nil
}

func (cp cardPair) compare_NoAppending() (bool, error) {
	return false, errors.New("No appendng")
}
func (cp cardPair) compare_NoMoving() (bool, error) {
	return false, errors.New("No moving")
}

func compare_Empty(dst *Pile, c *Card) (bool, error) {
	if dst.Label() != "" {
		if dst.Label() == "x" || dst.Label() == "X" {
			return false, errors.New("Cannot move cards to that empty pile")
		}
		ord := util.OrdinalToShortString(c.Ordinal())
		if ord != dst.Label() {
			return false, fmt.Errorf("Can only accept %s, not %s", util.ShortOrdinalToLongOrdinal(dst.Label()), util.ShortOrdinalToLongOrdinal(ord))
		}
	}
	return true, nil
}

// little library of simple compares

func (cp cardPair) compare_Up() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()-1 {
		return true, nil
	}
	return false, errors.New("Cards must be in ascending sequence")
}

func (cp cardPair) compare_UpWrap() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()-1 {
		return true, nil
	}
	if cp.c1.Ordinal() == 13 && cp.c2.Ordinal() == 1 {
		return true, nil // Ace on King
	}
	return false, errors.New("Cards must go up in rank (Aces on Kings allowed)")
}

func (cp cardPair) compare_Down() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()+1 {
		return true, nil
	}
	return false, errors.New("Cards must be in descending sequence")
}

func (cp cardPair) compare_DownWrap() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()+1 {
		return true, nil
	}
	if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	return false, errors.New("Cards must be in descending sequence (Kings on Aces allowed)")
}

func (cp cardPair) compare_UpOrDown() (bool, error) {
	if !(cp.c1.Ordinal()+1 == cp.c2.Ordinal() || cp.c1.Ordinal() == cp.c2.Ordinal()+1) {
		return false, errors.New("Cards must be in ascending or descending sequence")
	}
	return true, nil
}

func (cp cardPair) compare_UpOrDownWrap() (bool, error) {
	if (cp.c1.Ordinal()+1 == cp.c2.Ordinal()) || (cp.c1.Ordinal() == cp.c2.Ordinal()+1) {
		return true, nil
	} else if cp.c1.Ordinal() == 13 && cp.c2.Ordinal() == 1 {
		return true, nil // Ace On King
	} else if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	} else {
		return false, errors.New("Cards must be in ascending or descending sequence")
	}
}

func (cp cardPair) compare_Color() (bool, error) {
	if cp.c1.Black() != cp.c2.Black() {
		return false, errors.New("Cards must be the same color")
	}
	return true, nil
}

func (cp cardPair) compare_AltColor() (bool, error) {
	if cp.c1.Black() == cp.c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	return true, nil
}

func (cp cardPair) compare_Suit() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	return true, nil
}

func (cp cardPair) compare_OtherSuit() (bool, error) {
	if cp.c1.Suit() == cp.c2.Suit() {
		return false, errors.New("Cards must not be the same suit")
	}
	return true, nil
}

// library of compare functions made from simple compares

func (cp cardPair) compare_DownColor() (bool, error) {
	ok, err := cp.compare_Color()
	if !ok {
		return ok, err
	}
	return cp.compare_Down()
}

func (cp cardPair) compare_DownAltColor() (bool, error) {
	ok, err := cp.compare_AltColor()
	if !ok {
		return ok, err
	}
	return cp.compare_Down()
}

// compare_DownColorWrap not used
// func (cp cardPair) compare_DownColorWrap() (bool, error) {
// 	ok, err := cp.compare_Color()
// 	if !ok {
// 		return ok, err
// 	}
// 	return cp.compare_DownWrap()
// }

func (cp cardPair) compare_DownAltColorWrap() (bool, error) {
	ok, err := cp.compare_AltColor()
	if !ok {
		return ok, err
	}
	return cp.compare_DownWrap()
}

// compare_UpAltColor not used
// func (cp cardPair) compare_UpAltColor() (bool, error) {
// 	ok, err := cp.compare_AltColor()
// 	if !ok {
// 		return ok, err
// 	}
// 	return cp.compare_Up()
// }

func (cp cardPair) compare_UpSuit() (bool, error) {
	ok, err := cp.compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.compare_Up()
}

func (cp cardPair) compare_DownSuit() (bool, error) {
	ok, err := cp.compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.compare_Down()
}

func (cp cardPair) compare_UpOrDownSuit() (bool, error) {
	ok, err := cp.compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.compare_UpOrDown()
}

func (cp cardPair) compare_UpOrDownSuitWrap() (bool, error) {
	ok, err := cp.compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.compare_UpOrDownWrap()
}

// compare_DownOtherSuit not used
func (cp cardPair) compare_DownOtherSuit() (bool, error) {
	ok, err := cp.compare_OtherSuit()
	if !ok {
		return ok, err
	}
	return cp.compare_Down()
}

func (cp cardPair) compare_UpSuitWrap() (bool, error) {
	ok, err := cp.compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.compare_UpWrap()
}

func (cp cardPair) compare_DownSuitWrap() (bool, error) {
	ok, err := cp.compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.compare_DownWrap()
}

// chainCall
//
// Call using CardPair method expressions
// eg chainCall(CardPair.Compare_UpOrDown, CardPair.Compare_Suit)
//
// TODO think of something else for unsortedPairs(*Pile)
func (cp cardPair) chainCall(fns ...func(cardPair) (bool, error)) (ok bool, err error) {
	for _, fn := range fns {
		if ok, err = fn(cp); err != nil {
			break
		}
	}
	return
}
