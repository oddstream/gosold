//go:build linux || windows

package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosold/dark"
	light "oddstream.games/gosold/light"
	"oddstream.games/gosold/util"
)

func main() {

	var DebugMode bool

	log.SetFlags(0)

	// pearl from the mudbank: don't have any flags that will overwrite ThePreferences
	flag.BoolVar(&DebugMode, "debug", false, "turn debug mode on")
	flag.BoolVar(&dark.NoLoad, "noload", false, "do not load saved game when starting")
	flag.BoolVar(&dark.NoSave, "nosave", false, "do not save game before exit")
	var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()

	if DebugMode {
		dark.DebugMode = true
		light.DebugMode = true
		for i, a := range os.Args {
			log.Println(i, a)
		}
	}

	if *cpuProfile != "" {
		var f *os.File
		var err error
		if f, err = os.Create(*cpuProfile); err != nil {
			log.Fatal(err)
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	ebiten.SetScreenClearedEveryFrame(false)

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

	g := light.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

	// we come here if the user closed the window with the x button
	// but we don't come here if ExitRequested has been set
	// (and Game.Update() returned an error)
	// which another thing I don't understand
	// log.Println("main exit")
	g.ExitGame()
}
