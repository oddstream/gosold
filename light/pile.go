package light

import (
	"image"
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosold/dark"
	"oddstream.games/gosold/schriftbank"
)

const (
	CARD_FACE_FAN_FACTOR_V = 3.7
	CARD_FACE_FAN_FACTOR_H = 4
	CARD_BACK_FAN_FACTOR   = 8
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
	slot      image.Point  // local copy for mirror baize
	fanType   dark.FanType // local copy for mirror baize
	img       *ebiten.Image
}

func newPile(b *baize, darkPile *dark.Pile) *pile {
	return &pile{baize: b,
		darkPile:  darkPile,
		slot:      darkPile.Slot(),
		fanType:   darkPile.FanType(),
		fanFactor: defaultFanFactor[darkPile.FanType()]}
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

func (p *pile) push(c *card) {
	p.cards = append(p.cards, c)
	c.pile = p
	if darkProne := p.baize.darkBaize.IsCardProne(c.id); darkProne != c.lightProne {
		c.lightProne = darkProne
		c.startFlip()
	}
	c.lerpTo(c.pos)
}

func (p *pile) makeTail(c *card) []*card {
	if c.pile != p {
		log.Panic("Pile.makeTail called with a card that is not of this pile")
	}
	if c == p.peek() {
		return []*card{c}
	}
	for i, pc := range p.cards {
		if pc == c {
			return p.cards[i:]
		}
	}
	log.Panic("Pile.makeTail made an empty tail")
	return nil
}

// copyCardsFromDark
func (p *pile) copyCardsFromDark() {
	p.cards = []*card{}
	for _, id := range p.darkPile.Cards() {
		if c, ok := p.baize.cardMap[id.PackSuitOrdinal()]; !ok {
			log.Panicf("Card [%s] not found in card map", id.String())
		} else {
			p.push(c)
		}
	}
}

func (p *pile) createPlaceholder() {
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

// SetBaizePos sets the position of this Pile in Baize coords,
// and also sets the auxillary waste pile fanned positions
func (p *pile) setBaizePos(pos image.Point) {
	p.pos = pos
	switch p.darkPile.FanType() {
	case dark.FAN_DOWN3:
		p.pos1.X = p.pos.X
		p.pos1.Y = p.pos.Y + int(float64(CardHeight)/CARD_FACE_FAN_FACTOR_V)
		p.pos2.X = p.pos.X
		p.pos2.Y = p.pos1.Y + int(float64(CardHeight)/CARD_FACE_FAN_FACTOR_V)
	case dark.FAN_LEFT3:
		p.pos1.X = p.pos.X - int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		p.pos1.Y = p.pos.Y
		p.pos2.X = p.pos1.X - int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		p.pos2.Y = p.pos.Y
	case dark.FAN_RIGHT3:
		p.pos1.X = p.pos.X + int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		p.pos1.Y = p.pos.Y
		p.pos2.X = p.pos1.X + int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		p.pos2.Y = p.pos.Y
	}
}

// func (p *pile) baizePos() image.Point {
// 	return p.pos
// }

func (p *pile) screenPos() image.Point {
	return p.pos.Add(p.baize.dragOffset)
}

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

// PosAfter returns the position of the next card
func (p *pile) posAfter(c *card) image.Point {
	if len(p.cards) == 0 {
		println("Panic! PosAfter called in impossible way")
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
		// println("zero pos in PosAfter", self.category)
		return pos
	}
	switch p.darkPile.FanType() {
	case dark.FAN_NONE:
		// nothing to do
	case dark.FAN_DOWN:
		if p.baize.darkBaize.IsCardProne(c.id) {
			pos.Y += int(float64(CardHeight) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.Y += int(float64(CardHeight) / p.fanFactor)
		}
	case dark.FAN_LEFT:
		if p.baize.darkBaize.IsCardProne(c.id) {
			pos.X -= int(float64(CardWidth) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.X -= int(float64(CardWidth) / p.fanFactor)
		}
	case dark.FAN_RIGHT:
		if p.baize.darkBaize.IsCardProne(c.id) {
			pos.X += int(float64(CardWidth) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.X += int(float64(CardWidth) / p.fanFactor)
		}
	case dark.FAN_DOWN3, dark.FAN_LEFT3, dark.FAN_RIGHT3:
		switch len(p.cards) {
		case 0:
			// nothing to do
		case 1:
			pos = p.pos1 // incoming card at slot 1
		case 2:
			pos = p.pos2 // incoming card at slot 2
		default:
			pos = p.pos2 // incoming card at slot 2
			// top card needs to transition from slot[2] to slot[1]
			i := len(p.cards) - 1
			p.cards[i].lerpTo(p.pos1)
			// mid card needs to transition from slot[1] to slot[0]
			// all other cards to slot[0]
			for i > 0 {
				i--
				p.cards[i].lerpTo(p.pos)
			}
		}
	}
	return pos
}

func (p *pile) refan() {
	// TODO trying set pos instead of transition
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
func (p *pile) applyToCards(fn func(*card)) {
	for _, c := range p.cards {
		fn(c)
	}
}

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

func (p *pile) draw(screen *ebiten.Image) {

	if p.img == nil || p.darkPile.Hidden() {
		return
	}

	op := &ebiten.DrawImageOptions{}
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
}
