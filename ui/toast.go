package ui

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosold/schriftbank"
	"oddstream.games/gosold/sound"
)

/*
 https://material.io/archive/guidelines/components/snackbars-toasts.html

	Single-line snackbar height: 48dp
	Multi-line snackbar height: 80dp
	Text: Roboto Regular 14sp
	Action button: Roboto Medium 14sp, all-caps text
	Default background fill: #323232 100%
*/

// Toast represents a simple popup label that disappears after a few seconds
type Toast struct {
	img       *ebiten.Image
	message   string
	ticksLeft int
}

// ToastManager manages the list of toasts so that many can appear on screen at once
type ToastManager struct {
	toasts []*Toast
}

// Toast creates a new toast message an adds it to the list of messages
func (u *UI) Toast(soundEffect string, message string) {

	// play the sound even if the toast is already displayed
	sound.Play(soundEffect)

	// if we are already displaying this message, reset ticksLeft and quit
	// otherwise you can fill the screen with "Nothing to undo"
	for _, t := range u.toastManager.toasts {
		if t.message == message {
			t.ticksLeft = int(ebiten.ActualTPS()) * 6
			return
		}
	}

	// make a dummy context just to measure the width of thr message str9ing
	dc := gg.NewContext(8, 8)
	dc.SetFontFace(schriftbank.RobotoRegular14)
	w, _ := dc.MeasureString(message)

	w += 48
	h := float64(ToolbarHeight) // ignore measured height, force height to be 48

	// now reuse the dummy context, making it the correct size
	dc = gg.NewContext(int(w), int(h))
	dc.SetColor(BackgroundColor)
	dc.DrawRectangle(0, 0, w, h)
	dc.Fill()
	dc.Stroke()

	dc.SetFontFace(schriftbank.RobotoRegular14)
	dc.SetColor(ForegroundColor)
	dc.DrawStringAnchored(message, w/2, h/2, 0.5, 0.4)
	dc.Stroke()

	t := &Toast{message: message}
	t.img = ebiten.NewImageFromImage(dc.Image())
	// pearl from the mudbank, can't use ebiten.CurrentTPS() here
	// because during welcome toasts it will return 0.0
	// println(ebiten.CurrentTPS())
	// t.ticksLeft = int(ebiten.CurrentTPS()) * (6 + len(u.toastManager.toasts))
	t.ticksLeft = 60 * (8 + len(u.toastManager.toasts))

	u.toastManager.Add(t)
}

func (u *UI) ToastError(message string) {
	u.Toast("Error", message)
}

func (u *UI) ToastInfo(message string) {
	u.Toast("Glass", message)
}

// Add a new toast to the list
func (tm *ToastManager) Add(t *Toast) {
	tm.toasts = append(tm.toasts, t) // push onto end of list
	// println("Added toast", t.message)
}

// func (tm *ToastManager) Layout(outsideWidth, outsideHeight int) {
// }

// Update the queue of toasts
func (tm *ToastManager) Update() {
	if len(tm.toasts) == 0 {
		return
	}
	for _, t := range tm.toasts {
		t.ticksLeft--
	}
	for len(tm.toasts) > 0 && tm.toasts[0].ticksLeft < 0 {
		// println("Removing toast", tm.toasts[0].message)
		tm.toasts = tm.toasts[1:] // delete the oldest
	}
}

// Draw the toasts
func (tm *ToastManager) Draw(screen *ebiten.Image) {

	if len(tm.toasts) == 0 {
		// ebitenutil.DebugPrint(screen, "No toasts")
		return
	}
	// sx, sy := screen.Size()
	sx := screen.Bounds().Dx()
	sy := screen.Bounds().Dy()
	var tx, ty int
	ty = sy - 10 - 24 // 10 padding, 24 height of statusbar
	// for _, t := range tm.toasts {
	for i := len(tm.toasts) - 1; i >= 0; i-- {
		t := tm.toasts[i]
		// w, h := t.img.Size()
		w := t.img.Bounds().Dx()
		h := t.img.Bounds().Dy()
		tx = (sx - w) / 2
		ty = ty - h - 10 // move y up ready for next toast
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tx), float64(ty))
		screen.DrawImage(t.img, op)
		// ebitenutil.DebugPrint(screen, t.message)
	}
}
