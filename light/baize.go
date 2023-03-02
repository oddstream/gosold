package light

import "github.com/hajimehoshi/ebiten/v2"

type baize struct {
}

func (b *baize) layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (b *baize) update() error {
	return nil
}

func (b *baize) draw(screen *ebiten.Image) {

	screen.Fill(ExtendedColors["BaizeGreen"])
}
