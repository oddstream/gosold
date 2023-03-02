package light

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
	baize *baize
}

func NewGame() *Game {
	g := &Game{}
	g.baize = &baize{}
	return g
}

// Layout implements ebiten.Game's Layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.baize.layout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	g.baize.update()
	return nil
}

// Draw draws the current game to the given screen.
// Draw will be called based on the refresh rate of the screen (FPS).
// https://ebitencookbook.vercel.app/blog
func (g *Game) Draw(screen *ebiten.Image) {
	g.baize.draw(screen)
	// g.UI.Draw(screen)
}
