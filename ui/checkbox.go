package ui

import (
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosold/schriftbank"
)

// Checkbox is a button that displays a single rune
type Checkbox struct {
	WidgetBase
	boolVarPtr *bool
	fnUpdate   func()
	text       string
}

func (w *Checkbox) createImg() *ebiten.Image {
	dc := gg.NewContext(w.width, w.height)

	var iconName string
	if *(w.boolVarPtr) {
		iconName = "check_box"
	} else {
		iconName = "check_box_outline_blank"
	}
	// same as NavItem
	img, ok := IconMap[iconName]
	if !ok || img == nil {
		log.Fatal(iconName, " not in icon map")
	}
	dc.SetColor(ForegroundColor)
	dc.DrawImage(img, 0, w.height/4)
	dc.SetFontFace(schriftbank.RobotoMedium24)
	dc.DrawString(w.text, float64(48), float64(w.height)*0.8)

	// uncomment this to show the area we expect the text to occupy
	// dc.DrawLine(0, float64(0), float64(w.width), float64(0))
	// dc.DrawLine(0, float64(w.height), float64(w.width), float64(w.height))
	// dc.DrawLine(0, float64(0), float64(w.width), float64(w.height))
	// dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// NewCheckbox creates a new Checkbox
func NewCheckbox(parent Containery, id string, text string, boolVarPtr *bool, fnUpdate func()) *Checkbox {
	width, _ := parent.Size()
	w := &Checkbox{
		WidgetBase: WidgetBase{parent: parent, id: id, img: nil, x: 0, y: 0, width: width, height: 48},
		text:       text, boolVarPtr: boolVarPtr, fnUpdate: fnUpdate}
	w.Activate()
	return w
}

// Activate tells the input we need notifications
func (w *Checkbox) Activate() {
	w.disabled = false
	w.img = w.createImg()
}

// Deactivate tells the input we no longer need notofications
func (w *Checkbox) Deactivate() {
	w.disabled = true
	w.img = w.createImg()
}

func (w *Checkbox) Tapped() {
	if w.disabled {
		return
	}
	if w.boolVarPtr != nil {
		*(w.boolVarPtr) = !*(w.boolVarPtr)
	}
	w.img = w.createImg()
	if w.fnUpdate != nil {
		w.fnUpdate()
	}
	cmdFn(Command{Command: "SaveSettings"})
}
