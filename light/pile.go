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
	img       *ebiten.Image
}

func newPile(baize *baize, darkPile *dark.Pile) *pile {
	return &pile{baize: baize}
}

func (p *pile) updateCardsAndLabel(dp *dark.Pile) {
	p.darkPile = dp
	p.cards = nil
	for _, dc := range dp.Cards() {
		id := dc.ID().PackSuitOrdinal() // ignore prone flag
		if c, ok := p.baize.cardMap[id]; !ok {
			log.Panicf("Card %s not found in card map", id.String())
		} else {
			c.pile = p
			p.cards = append(p.cards, c)
		}
		// lerp card to it's pos
		// flip up or down
	}
	p.createPlaceholder()
}

func (p *pile) createPlaceholder() {
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
		if p.darkPile.Label() != "" {
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

		var label rune
		if p.baize.darkBaize.Recycles() == 0 {
			label = NORECYCLE_RUNE
		} else {
			label = RECYCLE_RUNE
		}
		dc.SetFontFace(schriftbank.CardSymbolHuge)
		dc.DrawStringAnchored(string(label), float64(CardWidth)*0.5, float64(CardHeight)*0.4, 0.5, 0.5)

		dc.Stroke()
		p.img = ebiten.NewImageFromImage(dc.Image())
	}
}
