package light

import (
	"errors"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosold/dark"
	"oddstream.games/gosold/sound"
)

var (
	// GosolVersionMajor is the integer version number
	GosoldVersionMajor int = 6
	// CsolVersionMinor is the integer version number
	GosoldVersionMinor int = 1
	// CSolVersionDate is the ISO 8601 date of bumping the version number
	GosolVersionDate string = "2023-03-03"
	// CardWidth of cards, start with a silly value to force a rescale/refan
	CardWidth int = 9
	// CardHeight of cards, start with a silly value to force a rescale/refan
	CardHeight int = 13
	// CardDiagonal float64 = 15.8
	// Card Corner Radius
	CardCornerRadius float64 = float64(CardWidth) / 10.0
	// PilePaddingX the gap left to the right of the pile
	PilePaddingX int = CardWidth / 10
	// PilePaddingY the gap left underneath each pile
	PilePaddingY int = CardHeight / 10
	// LeftMargin the gap between the left of the screen and the first pile
	LeftMargin int = (CardWidth / 2) + PilePaddingX
	// TopMargin the gap between top pile and top of baize
	TopMargin int = 48 /*ui.ToolbarHeight*/ + CardHeight/3
	// CardFaceImageLibrary
	// thirteen suitless cards,
	// one entry for each face card (4 suits * 13 cards),
	// suits are 1-indexed (eg club == 1) so image to be used for a card is (suit * 13) + (ord - 1).
	// can use (ord - 1) as an index to get suitless card
	TheCardFaceImageLibrary [13 * 5]*ebiten.Image
	// CardBackImage applies to all cards so is kept globally as an optimization
	CardBackImage *ebiten.Image
	// MovableCardBackImage applies to all cards so is kept globally as an optimization
	MovableCardBackImage *ebiten.Image
	// CardShadowImage applies to all cards so is kept globally as an optimization
	CardShadowImage *ebiten.Image
	// ExitRequested is set when user has had enough
	ExitRequested bool = false
)

type Game struct {
	darker       dark.Darker
	baize        *baize
	settings     *Settings
	commandTable map[ebiten.Key]func()
}

func NewGame() *Game {
	g := &Game{darker: dark.NewDark()}
	g.settings = NewSettings()
	g.settings.load()
	if g.settings.Mute {
		sound.SetVolume(0.0)
	} else {
		sound.SetVolume(g.settings.Volume)
	}
	g.baize = newBaize(g)
	g.baize.startGame(g.settings.Variant)

	g.commandTable = map[ebiten.Key]func(){
		ebiten.KeyC: func() { g.baize.collect() },
		ebiten.KeyN: func() { g.baize.newDeal() },
		ebiten.KeyR: func() { g.baize.restartDeal() },
		ebiten.KeyU: func() { g.baize.undo() },
		ebiten.KeyB: func() { g.baize.savePosition() },
		ebiten.KeyL: func() { g.baize.loadPosition() },
		ebiten.KeyS: func() { g.baize.savePosition() },
		ebiten.KeyX: func() { ExitRequested = true },
	}

	// TODO toast version bump
	return g
}

func (g *Game) execute(cmd any) {
	switch v := cmd.(type) {
	case ebiten.Key:
		if fn, ok := g.commandTable[v]; ok {
			fn()
		}
	default:
		log.Panicf("Game.execute unknown command type %v", cmd)
	}
}

// Layout implements ebiten.Game's Layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.baize.layout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	g.baize.update()
	// g.ui.update()
	if ExitRequested {
		g.baize.darkBaize.Save()
		g.settings.save()
		return errors.New("exit requested")
	}
	return nil
}

// Draw draws the current game to the given screen.
// Draw will be called based on the refresh rate of the screen (FPS).
// https://ebitencookbook.vercel.app/blog
func (g *Game) Draw(screen *ebiten.Image) {
	g.baize.draw(screen)
	// g.ui.draw(screen)
}
