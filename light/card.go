package light

import (
	"image"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosold/cardid"
	"oddstream.games/gosold/util"
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
	lerpingFlag   bool

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

	// card rotating things
	actualDeg, targetDeg int
}

func (c *card) prone() bool {
	return c.id.Prone()
}

func (c *card) setProne(prone bool) {
	c.id = c.id.SetProne(prone)
}

// baizePos returns the x,y baize coords of this card
func (c *card) baizePos() image.Point {
	return c.pos
}

// setBaizePos sets the position of the Card
func (c *card) setBaizePos(pos image.Point) {
	c.lerpingFlag = false
	c.pos = pos
}

// baizeRect gives the x,y baize coords of the card's top left and bottom right corners
func (c *card) baizeRect() image.Rectangle {
	var r image.Rectangle
	r.Min = c.pos
	r.Max = r.Min.Add(image.Point{CardWidth, CardHeight})
	return r
}

// screenRect gives the x,y screen coords of the card's top left and bottom right corners
func (c *card) screenRect() image.Rectangle {
	var r image.Rectangle = c.baizeRect()
	r.Min = r.Min.Add(c.pile.baize.dragOffset)
	r.Max = r.Max.Add(c.pile.baize.dragOffset)
	return r
}

// lerpTo starts the transition of this Card to pos
func (c *card) lerpTo(dst image.Point) {

	// if c.spinning() {
	// 	return
	// }

	if dst.Eq(c.pos) {
		c.lerpingFlag = false
		return // we are already here
	}

	if c.lerpingFlag && dst.Eq(c.dst) {
		return // repeat request to lerp to dst
	}

	c.lerpingFlag = true
	c.src = c.pos
	c.dst = dst
	// refanning waste cards can flutter with slow AniSpeed, so go faster if not far to go
	dist := util.Distance(c.src, c.dst)
	if dist < float64(CardWidth) {
		c.aniSpeed = c.pile.baize.game.settings.AniSpeed / 2.0
	} else {
		c.aniSpeed = c.pile.baize.game.settings.AniSpeed
	}
	c.lerpStartTime = time.Now()
}

// startDrag informs card that it is being dragged
func (c *card) startDrag() {
	if c.lerping() {
		// set the drag origin to the be transition destination,
		// so that cancelling this drag will return the card
		// to where it thought it was going
		// doing this will be trapped by Baize, so this is belt-n-braces
		c.dragStart = c.dst
	} else {
		c.dragStart = c.pos
	}
	c.beingDragged = true
	// println("start drag", c.ID.String(), "start", c.dragStartX, c.dragStartY)
}

// dragBy repositions the card by the distance it has been dragged
func (c *card) dragBy(dx, dy int) {
	// println("Card.DragBy(", c.dragStartX+dx-c.baizeX, c.dragStartY+dy-c.baizeY, ")")
	c.setBaizePos(c.dragStart.Add(image.Point{dx, dy}))
}

// DragStartPosition returns the x,y screen coords of this card before dragging started
// func (c *Card) DragStartPosition() (int, int) {
// return c.dragStartX, c.dragStartY
// }

// stopDrag informs card that it is no longer being dragged
func (c *card) stopDrag() {
	c.beingDragged = false
	// println("stop drag", c.ID.String())
}

// cancelDrag informs card that it is no longer being dragged
func (c *card) cancelDrag() {
	c.beingDragged = false
	// println("cancel drag", c.ID.String(), "start", c.dragStartX, c.dragStartY, "screen", c.screenX, c.screenY)
	c.lerpTo(c.dragStart)
}

// wasDragged returns true of this card has been dragged
func (c *card) wasDragged() bool {
	return !c.pos.Eq(c.dragStart)
}

func (c *card) startFlip() {
	c.flipWidth = 1.0    // card starts full width
	c.flipDirection = -1 // start by making card narrower
	c.flipStartTime = time.Now()
}

// startSpinning tells the card to start spinning
func (c *card) startSpinning() {
	c.directionX = rand.Intn(9) - 4
	c.directionY = rand.Intn(9) - 3 // favor falling downwards
	c.spin = rand.Float64() - 0.5
	// delay start of spinning to allow cards to be seen to go/finish their trip to foundations
	// https://stackoverflow.com/questions/67726230/creating-a-time-duration-from-float64-seconds
	d := time.Duration(c.pile.baize.game.settings.AniSpeed * float64(time.Second))
	d *= 2.0 // pause for admiration
	c.spinStartAfter = time.Now().Add(d)
}

// stopSpinning tells the card to stop spinning and return to it's upright state
func (c *card) stopSpinning() {
	c.directionX, c.directionY = 0, 0
	c.angle, c.spin = 0, 0
	// card may have spun off-screen slightly, and be -ve, which confuses Smoothstep
	c.pos = c.pile.pos
}

func (c *card) static() bool {
	return !c.lerpingFlag && !c.beingDragged && c.flipDirection == 0
}

// Spinning returns true if this card is spinning
func (c *card) spinning() bool {
	return c.spin != 0.0
}

// Lerping returns true if this card is lerping
func (c *card) lerping() bool {
	return c.lerpingFlag
}

// Dragging returns true if this card is being dragged
func (c *card) dragging() bool {
	return c.beingDragged
}

// Flipping returns true if this card is flipping
func (c *card) flipping() bool {
	return c.flipDirection != 0 // will be -1 or +1 if flipping
}

func (c *card) update() {
	if c.spinning() {
		if time.Now().After(c.spinStartAfter) {
			c.lerpingFlag = false
			c.pos.X += c.directionX
			c.pos.Y += c.directionY
			// pearl from the mudbank:
			// cannot flip card here (or anytime while spinning)
			// because Baize.Complete() will fail (and record a lost game)
			// because UnsortedPairs will "fail" because some cards will be face down
			// so do not call c.Flip() here
			c.angle += c.spin
			if c.angle > 360 {
				c.angle -= 360
			} else if c.angle < 0 {
				c.angle += 360
			}
		}
	}

	if c.lerping() {
		if !c.pos.Eq(c.dst) {
			secs := time.Since(c.lerpStartTime).Seconds()
			// secs will start at nearly zero, and rise to about the value of AniSpeed,
			// because AniSpeed is the number of seconds the card will take to transition.
			// with AniSpeed at 0.75, this happens (for example) 45 times (we are at @ 60Hz)
			var t float64 = secs / c.aniSpeed
			// with small values of AniSpeed, t can go above 1.0
			// which is bad: cards appear to fly away, never to be seen again
			// Smoothstep will correct this
			// if c.Ordinal() == 1 && c.Suit() == 1 {
			// 	log.Printf("%v\t0.25=%v\t0.5=%v\t0.75=%v", ts, ts/0.25, ts/0.5, ts/0.75)
			// }
			c.pos.X = int(util.Smoothstep(float64(c.src.X), float64(c.dst.X), t))
			c.pos.Y = int(util.Smoothstep(float64(c.src.Y), float64(c.dst.Y), t))
		} else {
			c.lerpingFlag = false
		}
	}

	if c.flipping() {
		// we need to flip faster than we lerp, because flipping happens in two stages
		t := time.Since(c.flipStartTime).Seconds() / (c.pile.baize.game.settings.AniSpeed / 2.0)
		if c.flipDirection < 0 {
			c.flipWidth = util.Lerp(1.0, 0.0, t)
			if c.flipWidth <= 0.0 {
				// reverse direction, make card wider
				c.flipDirection = 1
				c.flipStartTime = time.Now()
			}
		} else if c.flipDirection > 0 {
			c.flipWidth = util.Lerp(0.0, 1.0, t)
			if c.flipWidth >= 1.0 {
				c.flipDirection = 0
				c.flipWidth = 1.0
			}
		}
	}

	if c.actualDeg != c.targetDeg {
		// caveat: assumes slot.Deg values are multiples of 5 (15, 30, 45, 60, 75, 90)
		// could put a lerp here
		if c.targetDeg < c.actualDeg {
			c.actualDeg -= 5
		} else {
			c.actualDeg += 5
		}
	}
}

func (c *card) draw(screen *ebiten.Image) {
	if c.pile.hidden() {
		return
	}
	op := &ebiten.DrawImageOptions{}

	var img *ebiten.Image
	// card prone has already been set to destination state
	if c.flipDirection < 0 {
		if c.prone() {
			// card is getting narrower, and it's going to show face down, but show face up
			img = TheCardFaceImageLibrary[(c.id.Suit()*13)+(c.id.Ordinal()-1)]
		} else {
			// card is getting narrower, and it's going to show face up, but show face down
			img = CardBackImage
		}
	} else {
		if c.prone() {
			img = CardBackImage
		} else {
			img = TheCardFaceImageLibrary[(c.id.Suit()*13)+(c.id.Ordinal()-1)]
		}
	}

	if c.flipping() {
		op.GeoM.Translate(float64(-CardWidth/2), 0)
		op.GeoM.Scale(c.flipWidth, 1.0)
		op.GeoM.Translate(float64(CardWidth/2), 0)
	}

	if c.spinning() {
		// do this before the baize position translate
		op.GeoM.Translate(float64(-CardWidth/2), float64(-CardHeight/2))
		op.GeoM.Rotate(c.angle * 3.1415926535 / 180.0)
		op.GeoM.Translate(float64(CardWidth/2), float64(CardHeight/2))

		// naughty to do this here instead of Update(), but Draw() knows the screen dimensions and Update() doesn't
		// w, h := screen.Size()
		w := screen.Bounds().Dx() - c.pile.baize.dragOffset.X
		h := screen.Bounds().Dy() - c.pile.baize.dragOffset.Y
		switch {
		case c.pos.X+CardWidth > w:
			c.directionX = -rand.Intn(5)
			c.spin = rand.Float64() - 0.5
		case c.pos.X < 0:
			c.directionX = rand.Intn(5)
			c.spin = rand.Float64() - 0.5
		case c.pos.Y > h+CardHeight:
			c.directionX = rand.Intn(5)
			c.pos.Y = -CardHeight
		case c.pos.Y < -CardHeight:
			c.directionY = rand.Intn(5)
		}
	}

	// if c.pile.slot.Deg != 0 {
	if c.actualDeg != 0 {
		// rotate *before* translate
		// card is rotated about top left corner
		// and looks horribly jaggy
		// and the hit rect for it is, of course, wrong
		op.GeoM.Translate(-float64(CardWidth)/2, -float64(CardHeight)/2)
		op.GeoM.Rotate(float64(c.actualDeg) * math.Pi / 180.0)
		op.GeoM.Translate(float64(CardWidth)/2, float64(CardHeight)/2)
	}

	op.GeoM.Translate(float64(c.pos.X+c.pile.baize.dragOffset.X), float64(c.pos.Y+c.pile.baize.dragOffset.Y))

	if !c.flipping() {
		if c.lerping() || c.dragging() {
			op.GeoM.Translate(4.0, 4.0)
			screen.DrawImage(CardShadowImage, op)
			op.GeoM.Translate(-4.0, -4.0)
		}
		// no longer "press" the card when dragging it
		// because this made tapping look a little messy
	}

	if c.pile.baize.game.settings.ShowMovableCards && !c.spinning() {
		if c.pile.darkPile.IsStock() {
			// card will be prone because Stock
			// nb this will color all the stock cards, not just the top card
			img = MovableCardBackImage
		} else {
			var weight int16 = c.pile.baize.darkBaize.CardTapWeight(c.id)
			if !c.flipping() && weight != 0 {
				var fudgeFactor float32 = 1.0 - (float32(weight) / 10)
				op.ColorScale.Scale(1.0, 1.0, fudgeFactor, 1.0)
				// switch weight {
				// case 1: // Cell
				// 	op.ColorScale.Scale(1.0, 1.0, 0.9, 1)
				// case 2: // Normal
				// 	op.ColorScale.Scale(1.0, 1.0, 0.8, 1)
				// case 3: // Open pile or turn up
				// 	op.ColorScale.Scale(1.0, 1.0, 0.7, 1)
				// case 4: // Suit match
				// 	op.ColorScale.Scale(1.0, 1.0, 0.6, 1)
				// case 5: // Discard or Foundation
				// 	op.ColorScale.Scale(1.0, 1.0, 0.5, 1)
				// }
			}
		}
	}

	if img != nil {
		screen.DrawImage(img, op)
		// if DebugMode && c.pile.baize.darkBaize.CardTapWeight(c.id) > 0 {
		// 	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", c.pile.baize.darkBaize.CardTapWeight(c.id)), c.pos.X+c.pile.baize.dragOffset.X, c.pos.Y+c.pile.baize.dragOffset.Y)
		// }
	}
}
