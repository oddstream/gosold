package ui

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosold/schriftbank"
)

// TextUrl is a widget that displays a clickable url
type TextUrl struct {
	WidgetBase
	text string
	url  string
}

func (w *TextUrl) createImg() *ebiten.Image {
	dc := gg.NewContext(w.width, w.height)

	dc.SetRGBA(0.5, 0.5, 1, 1)
	// nota bene - text is drawn with y as a baseline
	dc.SetFontFace(schriftbank.RobotoRegular14)
	dc.DrawString(w.text, 0, float64(w.height-8)) // move up a little to stop descenders being clipped
	// uncomment this line to visualize text box
	// dc.DrawLine(0, 0, float64(w.width), float64(w.height))
	// dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

func (w *TextUrl) calcHeights() {
	dc := gg.NewContext(w.width, 48)
	dc.SetFontFace(schriftbank.RobotoRegular14)
	w.height = 24
}

// NewTextUrl creates a new TextUrl widget
func NewTextUrl(parent Containery, id string, text string) *TextUrl {
	width, _ := parent.Size()
	// widget x, y will be set by LayoutWidgets
	// widget height will be set when wordwrapping in createImg
	w := &TextUrl{
		WidgetBase: WidgetBase{parent: parent, id: id, img: nil, width: width},
		text:       "Wikipedia",
		url:        text}
	w.calcHeights()
	w.Activate()
	return w
}

// Activate tells the input we need notifications
func (w *TextUrl) Activate() {
	w.disabled = false
	w.img = w.createImg()
	// w.input.Add(w)
}

// Deactivate tells the input we no longer need notofications
func (w *TextUrl) Deactivate() {
	w.disabled = true
	w.img = w.createImg()
	// w.input.Remove(w)
}

func (w *TextUrl) Tapped() {
	if w.disabled {
		return
	}
	OpenBrowserWindow(w.url)
}
