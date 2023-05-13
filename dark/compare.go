package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"

	"oddstream.games/gosold/util"
)

type dyad struct {
	c1, c2 *Card
}

type dyadCmpFunc func(dyad) (bool, error)

func tailConformant(tail []*Card, fn dyadCmpFunc) (bool, error) {
	for i := 1; i < len(tail); i++ {
		if ok, err := fn(dyad{tail[i-1], tail[i]}); !ok {
			return false, err
		}
	}
	return true, nil
}

func (dyad) compare_NoAppending() (bool, error) {
	return false, errors.New("No appending")
}
func (dyad) compare_NoMoving() (bool, error) {
	return false, errors.New("No moving")
}

func compare_Empty(dst *Pile, tail []*Card) (bool, error) {
	if dst.Label() != "" {
		if dst.Label() == "X" {
			return false, errors.New("Cannot move cards to that empty pile")
		}
		ord := util.OrdinalToShortString(tail[0].Ordinal())
		if ord != dst.Label() {
			return false, fmt.Errorf("Can only accept %s, not %s", util.ShortOrdinalToLongOrdinal(dst.Label()), util.ShortOrdinalToLongOrdinal(ord))
		}
	}
	return true, nil
}

// little library of simple compares

func (dyad) compare_Any() (bool, error) {
	return true, nil
}

func (dy dyad) compare_Up() (bool, error) {
	if dy.c1.Ordinal() == dy.c2.Ordinal()-1 {
		return true, nil
	}
	return false, errors.New("Cards must be in ascending sequence")
}

func (dy dyad) compare_UpWrap() (bool, error) {
	if dy.c1.Ordinal() == dy.c2.Ordinal()-1 {
		return true, nil
	}
	if dy.c1.Ordinal() == 13 && dy.c2.Ordinal() == 1 {
		return true, nil // Ace on King
	}
	return false, errors.New("Cards must go up in rank (Aces on Kings allowed)")
}

func (dy dyad) compare_Down() (bool, error) {
	if dy.c1.Ordinal() == dy.c2.Ordinal()+1 {
		return true, nil
	}
	return false, errors.New("Cards must be in descending sequence")
}

func (dy dyad) compare_DownWrap() (bool, error) {
	if dy.c1.Ordinal() == dy.c2.Ordinal()+1 {
		return true, nil
	}
	if dy.c1.Ordinal() == 1 && dy.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	return false, errors.New("Cards must be in descending sequence (Kings on Aces allowed)")
}

func (dy dyad) compare_UpOrDown() (bool, error) {
	if !(dy.c1.Ordinal()+1 == dy.c2.Ordinal() || dy.c1.Ordinal() == dy.c2.Ordinal()+1) {
		return false, errors.New("Cards must be in ascending or descending sequence")
	}
	return true, nil
}

func (dy dyad) compare_UpOrDownWrap() (bool, error) {
	if (dy.c1.Ordinal()+1 == dy.c2.Ordinal()) || (dy.c1.Ordinal() == dy.c2.Ordinal()+1) {
		return true, nil
	} else if dy.c1.Ordinal() == 13 && dy.c2.Ordinal() == 1 {
		return true, nil // Ace On King
	} else if dy.c1.Ordinal() == 1 && dy.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	} else {
		return false, errors.New("Cards must be in ascending or descending sequence")
	}
}

func (dy dyad) compare_Color() (bool, error) {
	if dy.c1.Black() != dy.c2.Black() {
		return false, errors.New("Cards must be the same color")
	}
	return true, nil
}

func (dy dyad) compare_AltColor() (bool, error) {
	if dy.c1.Black() == dy.c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	return true, nil
}

func (dy dyad) compare_Suit() (bool, error) {
	if dy.c1.Suit() != dy.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	return true, nil
}

func (dy dyad) compare_OtherSuit() (bool, error) {
	if dy.c1.Suit() == dy.c2.Suit() {
		return false, errors.New("Cards must not be the same suit")
	}
	return true, nil
}

// library of compare functions made from simple compares

func (dy dyad) compare_DownColor() (bool, error) {
	// ok, err := dy.compare_Color()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_Down()
	return dy.chainCall(dyad.compare_Color, dyad.compare_Down)
}

func (dy dyad) compare_DownAltColor() (bool, error) {
	// ok, err := dy.compare_AltColor()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_Down()
	return dy.chainCall(dyad.compare_AltColor, dyad.compare_Down)
}

// compare_DownColorWrap not used
// func (cp cardPair) compare_DownColorWrap() (bool, error) {
// 	ok, err := cp.compare_Color()
// 	if !ok {
// 		return ok, err
// 	}
// 	return cp.compare_DownWrap()
// }

func (dy dyad) compare_DownAltColorWrap() (bool, error) {
	// ok, err := dy.compare_AltColor()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_DownWrap()
	return dy.chainCall(dyad.compare_AltColor, dyad.compare_DownWrap)
}

func (dy dyad) compare_UpAltColor() (bool, error) {
	// ok, err := dy.compare_AltColor()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_Up()
	return dy.chainCall(dyad.compare_AltColor, dyad.compare_Up)
}

func (dy dyad) compare_UpSuit() (bool, error) {
	// ok, err := dy.compare_Suit()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_Up()
	return dy.chainCall(dyad.compare_Suit, dyad.compare_Up)
}

func (dy dyad) compare_DownSuit() (bool, error) {
	// ok, err := dy.compare_Suit()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_Down()
	return dy.chainCall(dyad.compare_Suit, dyad.compare_Down)
}

func (dy dyad) compare_UpOrDownSuit() (bool, error) {
	// ok, err := dy.compare_Suit()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_UpOrDown()
	return dy.chainCall(dyad.compare_Suit, dyad.compare_UpOrDown)
}

func (dy dyad) compare_UpOrDownSuitWrap() (bool, error) {
	// ok, err := dy.compare_Suit()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_UpOrDownWrap()
	return dy.chainCall(dyad.compare_Suit, dyad.compare_UpOrDownWrap)
}

// compare_DownOtherSuit not used
func (dy dyad) compare_DownOtherSuit() (bool, error) {
	// ok, err := dy.compare_OtherSuit()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_Down()
	return dy.chainCall(dyad.compare_OtherSuit, dyad.compare_Down)
}

func (dy dyad) compare_UpSuitWrap() (bool, error) {
	// ok, err := dy.compare_Suit()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_UpWrap()
	return dy.chainCall(dyad.compare_Suit, dyad.compare_UpWrap)
}

func (dy dyad) compare_DownSuitWrap() (bool, error) {
	// ok, err := dy.compare_Suit()
	// if !ok {
	// 	return ok, err
	// }
	// return dy.compare_DownWrap()
	return dy.chainCall(dyad.compare_Suit, dyad.compare_DownWrap)
}

// chainCall
//
// Call using dyad method expressions
// eg chainCall(dyad.compare_UpOrDown, dyad.compare_Suit)
func (dy dyad) chainCall(fns ...func(dyad) (bool, error)) (ok bool, err error) {
	for _, fn := range fns {
		if ok, err = fn(dy); err != nil {
			break
		}
	}
	return
}
