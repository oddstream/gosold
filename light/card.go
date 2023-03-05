package light

import (
	"image"
	"time"

	"oddstream.games/gosold/cardid"
	"oddstream.games/gosold/dark"
)

type card struct {
	id   cardid.CardID
	pile *pile

	// pos of card on light.baize
	pos image.Point

	// lerping things
	src           image.Point // lerp origin
	dst           image.Point // lerp destination
	aniSpeed      float64
	lerpStartTime time.Time
	lerping       bool

	// dragging things
	dragStart    image.Point // starting point for dragging
	beingDragged bool        // true if this card is being dragged, or is in a dragged tail

	// flipping things
	flipWidth     float64 // scale of the card width while flipping
	flipDirection int
	flipStartTime time.Time

	// spinning things
	directionX, directionY int     // direction vector when card is spinning
	angle, spin            float64 // current angle and spin when card is spinning
	spinStartAfter         time.Time
}

func newCard(darkCard *dark.Card) *card {
	return &card{id: darkCard.ID()} // pos will start at 0,0
}
