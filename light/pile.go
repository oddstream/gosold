package light

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"oddstream.games/gosold/dark"
	"oddstream.games/gosold/schriftbank"
	"oddstream.games/gosold/util"
)

const (
	CARD_FACE_FAN_FACTOR_V = 1.0 / 3.7
	CARD_FACE_FAN_FACTOR_H = 1.0 / 4.0
	CARD_BACK_FAN_FACTOR   = 1.0 / 8.0
	MIN_FACE_FAN_FACTOR    = 1.0 / 6.0
)

var defaultFanFactor [7]float64 = [7]float64{
	1.0,                    // FAN_NONE
	CARD_FACE_FAN_FACTOR_V, // FAN_DOWN
	CARD_FACE_FAN_FACTOR_H, // FAN_LEFT,
	CARD_FACE_FAN_FACTOR_H, // FAN_RIGHT,
	CARD_FACE_FAN_FACTOR_V, // FAN_DOWN3,
	CARD_FACE_FAN_FACTOR_H, // FAN_LEFT3,
	CARD_FACE_FAN_FACTOR_H, // FAN_RIGHT3,
}

const (
	// https://en.wikipedia.org/wiki/Miscellaneous_Symbols
	RECYCLE_RUNE   = rune(0x267B)
	NORECYCLE_RUNE = rune(0x2613)
)

type pile struct {
	baize     *baize
	darkPile  *dark.Pile // use to get Slot(), FanType()
	cards     []*card
	pos       image.Point // actual position on baize
	pos1      image.Point // waste pos #1
	pos2      image.Point // waste pos #1
	fanFactor float64
	slot      dark.PileSlot // local copy for mirror baize
	fanType   dark.FanType  // local copy for mirror baize
	boundary  *pile
	box       image.Rectangle
	img       *ebiten.Image // placeholder
}

func newPile(b *baize, darkPile *dark.Pile) *pile {
	var p *pile = &pile{baize: b,
		darkPile:  darkPile,
		slot:      darkPile.Slot(),
		fanType:   darkPile.FanType(),
		fanFactor: defaultFanFactor[darkPile.FanType()]}
	return p
}

// func (p *pile) reset() {
// 	p.cards = p.cards[:0]
// 	p.fanFactor = defaultFanFactor[p.fanType]
// }

func (p *pile) peek() *card {
	if len(p.cards) == 0 {
		return nil
	}
	return p.cards[len(p.cards)-1]
}

// push a card onto this pile, flipping the card if light and dark prone
// flags do not match
func (p *pile) push(c *card) {
	p.cards = append(p.cards, c)
	c.pile = p
	if darkProne := p.baize.darkBaize.IsCardProne(c.id); darkProne != c.prone() {
		c.setProne(darkProne)
		c.startFlip()
	}
	c.targetDeg = p.slot.Deg
	c.lerpTo(c.pos)
}

// makeTail returns a slice of cards from c downwards
func (p *pile) makeTail(c *card) []*card {
	if DebugMode && c.pile != p {
		log.Panic("pile.makeTail called with a card that is not of this pile")
	}
	if c == p.peek() {
		return []*card{c}
	}
	for i, pc := range p.cards {
		if pc == c {
			return p.cards[i:]
		}
	}
	log.Panic("pile.makeTail made an empty tail")
	return nil
}

// copyCardsFromDark resets this pile, copies all the cards from dark to light,
// and pushes them onto this pile
func (p *pile) copyCardsFromDark() {
	p.cards = []*card{}
	for _, id := range p.darkPile.Cards() {
		if c, ok := p.baize.cardMap[id.PackSuitOrdinal()]; !ok {
			log.Panicf("%s not found in card map", id.String())
		} else {
			p.push(c)
		}
	}
	// fanning is done once, later
}

// hidden returns true if this pile is not displayed on screen
func (p *pile) hidden() bool {
	return p.slot.X < 0 || p.slot.Y < 0
}

// createPlaceHolder for this pile, depending on category. Sets pile.img field.
func (p *pile) createPlaceholder() {
	if p.hidden() || p.darkPile.Label() == "X" {
		return
	}
	// BEWARE the card fonts may not yet be loaded
	switch p.darkPile.Category() {
	case "Cell":
		// basic empty rounded rect
		dc := gg.NewContext(CardWidth, CardHeight)
		dc.SetColor(color.NRGBA{255, 255, 255, 31})
		dc.SetLineWidth(2)
		// draw the RoundedRect entirely INSIDE the context
		dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
		dc.Stroke()
		p.img = ebiten.NewImageFromImage(dc.Image())
	case "Discard":
		// filled in rounded rect
		dc := gg.NewContext(CardWidth, CardHeight)
		dc.SetColor(color.NRGBA{255, 255, 255, 31})
		dc.SetLineWidth(2)
		// draw the RoundedRect entirely INSIDE the context
		dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
		dc.Fill() // difference for this subpile
		dc.Stroke()
		p.img = ebiten.NewImageFromImage(dc.Image())
	case "Foundation", "Tableau":
		// basic empty rounded rect, with label
		dc := gg.NewContext(CardWidth, CardHeight)
		dc.SetColor(color.NRGBA{255, 255, 255, 31})
		dc.SetLineWidth(2)
		// draw the RoundedRect entirely INSIDE the context
		dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
		if p.darkPile.Label() != "" && schriftbank.CardOrdinalLarge != nil {
			dc.SetFontFace(schriftbank.CardOrdinalLarge)
			dc.DrawStringAnchored(p.darkPile.Label(), float64(CardWidth)*0.5, float64(CardHeight)*0.4, 0.5, 0.5)
		}
		dc.Stroke()
		p.img = ebiten.NewImageFromImage(dc.Image())
	case "Reserve", "Waste":
		// nothing, p.img stays nil
	case "Stock":
		// basic empty rounded rect, with recycle rune
		dc := gg.NewContext(CardWidth, CardHeight)
		dc.SetColor(color.NRGBA{255, 255, 255, 31})
		dc.SetLineWidth(2)
		// draw the RoundedRect entirely INSIDE the context
		dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)

		// farted around trying to use icons for this
		// but they were 48x48 and got fuzzy when scaled
		// and were stubbornly white

		if schriftbank.CardSymbolHuge != nil {
			var label rune // gimmie a ternery operator
			if p.baize.darkBaize.Recycles() == 0 {
				label = NORECYCLE_RUNE
			} else {
				label = RECYCLE_RUNE
			}
			dc.SetFontFace(schriftbank.CardSymbolHuge)
			dc.DrawStringAnchored(string(label), float64(CardWidth)*0.5, float64(CardHeight)*0.4, 0.5, 0.5)
		}

		dc.Stroke()
		p.img = ebiten.NewImageFromImage(dc.Image())
	}
}

// setBaizePos sets the position of this Pile in Baize coords,
// and also sets the auxillary waste pile fanned positions
func (p *pile) setBaizePos(pos image.Point) {
	p.pos = pos
	switch p.darkPile.FanType() {
	case dark.FAN_DOWN3:
		p.pos1.X = p.pos.X
		p.pos1.Y = p.pos.Y + int(float64(CardHeight)*CARD_FACE_FAN_FACTOR_V)
		p.pos2.X = p.pos.X
		p.pos2.Y = p.pos1.Y + int(float64(CardHeight)*CARD_FACE_FAN_FACTOR_V)
	case dark.FAN_LEFT3:
		p.pos1.X = p.pos.X - int(float64(CardWidth)*CARD_FACE_FAN_FACTOR_H)
		p.pos1.Y = p.pos.Y
		p.pos2.X = p.pos1.X - int(float64(CardWidth)*CARD_FACE_FAN_FACTOR_H)
		p.pos2.Y = p.pos.Y
	case dark.FAN_RIGHT3:
		p.pos1.X = p.pos.X + int(float64(CardWidth)*CARD_FACE_FAN_FACTOR_H)
		p.pos1.Y = p.pos.Y
		p.pos2.X = p.pos1.X + int(float64(CardWidth)*CARD_FACE_FAN_FACTOR_H)
		p.pos2.Y = p.pos.Y
	}
}

// func (p *pile) baizePos() image.Point {
// 	return p.pos
// }

// func (p *pile) screenPos() image.Point {
// 	return p.pos.Add(p.baize.dragOffset)
// }

func (p *pile) baizeRect() image.Rectangle {
	var r image.Rectangle
	r.Min = p.pos
	r.Max = r.Min.Add(image.Point{CardWidth, CardHeight})
	return r
}

func (p *pile) screenRect() image.Rectangle {
	var r image.Rectangle = p.baizeRect()
	r.Min = r.Min.Add(p.baize.dragOffset)
	r.Max = r.Max.Add(p.baize.dragOffset)
	return r
}

func (p *pile) fannedBaizeRect() image.Rectangle {
	var r image.Rectangle = p.baizeRect()
	if len(p.cards) > 1 {
		var c *card = p.peek()
		// if c.Dragging() {
		// 	return r
		// }
		var cPos = c.baizePos()
		switch p.darkPile.FanType() {
		case dark.FAN_NONE:
			// do nothing
		case dark.FAN_RIGHT, dark.FAN_RIGHT3:
			r.Max.X = cPos.X + CardWidth
		case dark.FAN_LEFT, dark.FAN_LEFT3:
			r.Max.X = cPos.X - CardWidth
		case dark.FAN_DOWN, dark.FAN_DOWN3:
			r.Max.Y = cPos.Y + CardHeight
		}
	}
	return r
}

// func (p *pile) fannedScreenRect() image.Rectangle {
// 	var r image.Rectangle = p.fannedBaizeRect()
// 	r.Min = r.Min.Add(p.baize.dragOffset)
// 	r.Max = r.Max.Add(p.baize.dragOffset)
// 	return r
// }

func (p *pile) calcBoundaryBox() {
	if p.boundary == nil {
		return
	}
	switch p.fanType {
	case dark.FAN_DOWN:
		p.box = image.Rectangle{
			Min: p.pos, // image.Point{p.pos.X, p.pos.Y},
			Max: image.Point{p.pos.X + CardWidth, p.boundary.pos.Y},
		}
	case dark.FAN_RIGHT:
		p.box = image.Rectangle{
			Min: p.pos, // image.Point{p.pos.X, p.pos.Y},
			Max: image.Point{p.boundary.pos.X, p.pos.Y + CardHeight},
		}
	case dark.FAN_LEFT:
		p.box = image.Rectangle{
			Min: image.Point{p.boundary.pos.X, p.pos.Y},
			Max: image.Point{p.pos.X + CardWidth, p.pos.Y + CardHeight},
		}
	}
}

func (p *pile) calcFaceFanFactor() {
	/*
		result = ((#cards - 1) * (cardheight * factor)) + cardheight
		r = (n-1) * (h * f) + h
		make factor the subject
		f = (r - h) / (h * (n-1))
		https://www.mymathtutors.com/algebra-tutors/adding-numerators/online-calculator---rearrange.html
	*/
	if p.boundary == nil || len(p.cards) < 2 {
		p.fanFactor = defaultFanFactor[p.fanType]
		return
	}
	// p.box will already be set, in screen coords (TODO check)
	var ff float64 = defaultFanFactor[p.fanType]
	switch p.fanType {
	case dark.FAN_DOWN:
		ff = float64(p.box.Dy()-CardHeight) / float64(CardHeight*(len(p.cards)-1))
		ff = util.Clamp(ff, MIN_FACE_FAN_FACTOR, CARD_FACE_FAN_FACTOR_V)
	case dark.FAN_RIGHT, dark.FAN_LEFT:
		ff = float64(p.box.Dx()-CardWidth) / float64(CardWidth*(len(p.cards)-1))
		ff = util.Clamp(ff, MIN_FACE_FAN_FACTOR, CARD_FACE_FAN_FACTOR_H)
	}
	p.fanFactor = ff
}

// posAfter returns the position of the next card after c
func (p *pile) posAfter(c *card) image.Point {
	if DebugMode && len(p.cards) == 0 {
		log.Panic("pile.posAfter called in impossible way")
		return p.pos
	}
	var pos image.Point
	if c.lerping() {
		pos = c.dst
	} else {
		pos = c.pos
	}
	if pos.X <= 0 && pos.Y <= 0 {
		// the card is still at 0,0 where it started life
		// and is yet to have pos calculated from the pile slot
		// println("zero pos in posAfter", self.category)
		return pos
	}
	switch p.darkPile.FanType() {
	case dark.FAN_NONE:
		// nothing to do
	case dark.FAN_DOWN:
		if c.prone() {
			pos.Y += int(float64(CardHeight) * float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.Y += int(float64(CardHeight) * p.fanFactor)
		}
	case dark.FAN_LEFT:
		if c.prone() {
			pos.X -= int(float64(CardWidth) * float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.X -= int(float64(CardWidth) * p.fanFactor)
		}
	case dark.FAN_RIGHT:
		if c.prone() {
			pos.X += int(float64(CardWidth) * float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.X += int(float64(CardWidth) * p.fanFactor)
		}
	case dark.FAN_DOWN3, dark.FAN_LEFT3, dark.FAN_RIGHT3:
		switch len(p.cards) {
		case 0:
			// nothing to do
		case 1:
			pos = p.pos1 // incoming card at position 1
		case 2:
			pos = p.pos2 // incoming card at position 2
		default:
			pos = p.pos2 // incoming card at position 2
			// top card needs to transition from position[2] to position[1]
			i := len(p.cards) - 1
			p.cards[i].lerpTo(p.pos1)
			// mid card needs to transition from position[1] to position[0]
			// all other cards to position[0]
			for i > 0 {
				i--
				p.cards[i].lerpTo(p.pos)
			}
		}
	}
	return pos
}

// refan repositions all the cards in this pile
func (p *pile) refan() {
	p.calcFaceFanFactor() // do this before using posAfter()
	var doFan3 bool = false
	switch p.darkPile.FanType() {
	case dark.FAN_NONE:
		for _, c := range p.cards {
			c.lerpTo(p.pos)
		}
	case dark.FAN_DOWN3, dark.FAN_LEFT3, dark.FAN_RIGHT3:
		for _, c := range p.cards {
			c.lerpTo(p.pos)
		}
		doFan3 = true
	case dark.FAN_DOWN, dark.FAN_LEFT, dark.FAN_RIGHT:
		var pos = p.pos
		var i = 0
		for _, c := range p.cards {
			c.lerpTo(pos)
			pos = p.posAfter(p.cards[i])
			i++
		}
	}

	if doFan3 {
		switch len(p.cards) {
		case 0:
		case 1:
			// nothing to do
		case 2:
			c := p.cards[1]
			c.lerpTo(p.pos1)
		default:
			i := len(p.cards)
			i--
			c := p.cards[i]
			c.lerpTo(p.pos2)
			i--
			c = p.cards[i]
			c.lerpTo(p.pos1)
		}
	}
}

// applyToCards applies a function to each card in the pile
// caller must use a method expression, eg (*Card).StartSpinning, yielding a function value
// with a regular first parameter taking the place of the receiver
// func (p *pile) applyToCards(fn func(*card)) {
// 	for _, c := range p.cards {
// 		fn(c)
// 	}
// }

func (p *pile) drawStaticCards(screen *ebiten.Image) {
	for _, c := range p.cards {
		if c.static() {
			c.draw(screen)
		}
	}
}

func (p *pile) drawAnimatingCards(screen *ebiten.Image) {
	for _, c := range p.cards {
		if c.lerping() || c.flipping() {
			c.draw(screen)
		}
	}
}

func (p *pile) drawDraggingCards(screen *ebiten.Image) {
	for _, c := range p.cards {
		if c.dragging() {
			c.draw(screen)
		}
	}
}

func (p *pile) update() {
	for _, card := range p.cards {
		card.update()
	}
}

// draw the pile's placeHolder
func (p *pile) draw(screen *ebiten.Image) {

	if p.img == nil {
		// hidden piles will not have a placeHolder
		return
	}

	op := &ebiten.DrawImageOptions{}

	if p.slot.Deg != 0 {
		op.GeoM.Translate(-float64(CardWidth)/2, -float64(CardHeight)/2)
		op.GeoM.Rotate(float64(p.slot.Deg) * math.Pi / 180.0)
		op.GeoM.Translate(float64(CardWidth)/2, float64(CardHeight)/2)
	}

	op.GeoM.Translate(float64(p.pos.X+p.baize.dragOffset.X), float64(p.pos.Y+p.baize.dragOffset.Y))
	// if self.target && len(self.cards) == 0 {
	// 	op.ColorScale.Scale(0.75, 0.75, 0.75, 1)
	// 	// op.GeoM.Translate(2, 2)
	// }

	if p.darkPile.IsStock() && p.baize.darkBaize.Recycles() > 0 {
		if pt := image.Pt(ebiten.CursorPosition()); pt.In(p.screenRect()) {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				op.GeoM.Translate(2, 2)
			}
		}
	}

	screen.DrawImage(p.img, op)

	if DebugMode && DrawBoxes && p.boundary != nil {
		var x float32 = float32(p.box.Min.X + p.baize.dragOffset.X) // effectively p.screenPos()
		var y float32 = float32(p.box.Min.Y + p.baize.dragOffset.Y)
		var width float32 = float32(p.box.Dx())
		var height float32 = float32(p.box.Dy())
		vector.DrawFilledRect(screen, x, y, width, height, color.NRGBA{255, 255, 255, 31}, false)
	}
}
