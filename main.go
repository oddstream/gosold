//go:build linux || windows

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	light "oddstream.games/gosold/light"
	"oddstream.games/gosold/util"
)

func main() {
	// ebiten panics if a window to maximize is not resizable
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if ebiten.IsWindowMaximized() || ebiten.IsWindowMinimized() {
		// GNOME (maybe) annoyingly keeps maximizing the window
		ebiten.RestoreWindow()
	}
	{
		x, y := ebiten.ScreenSizeInFullscreen()
		n := util.Max(x, y)
		ebiten.SetWindowSize(n/2, n/2)
	}
	ebiten.SetWindowIcon(light.WindowIcons())
	ebiten.SetWindowTitle("Go Solitaire")

	g := &light.Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
