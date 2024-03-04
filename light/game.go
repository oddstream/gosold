package light

import (
	"errors"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosold/dark"
	"oddstream.games/gosold/sound"
	"oddstream.games/gosold/ui"
)

var (
	DebugMode bool
	// GosolVersionMajor is the integer version number
	GosoldVersionMajor int = 6
	// CsolVersionMinor is the integer version number
	GosoldVersionMinor int = 5
	// CSolVersionDate is the ISO 8601 date of bumping the version number
	GosoldVersionDate string = "2024-03-04"
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
	TopMargin int = ui.ToolbarHeight + CardHeight/3
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
	ExitRequested bool
	// DrawBoxes displays the pile boundary boxes
	DrawBoxes bool
)

type Game struct {
	darker       dark.Darker
	ui           *ui.UI
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
	g.ui = ui.New(g.execute)
	g.baize = newBaize(g)
	g.baize.startGame(g.settings.Variant)
	g.baize.darkBaize.Load()

	g.commandTable = map[ebiten.Key]func(){
		ebiten.KeyN: func() { g.baize.newDeal() },
		ebiten.KeyR: func() { g.baize.restartDeal() },
		ebiten.KeyF: func() { g.ui.ShowVariantPickerEx(g.darker.ListVariantGroups(), "ShowVariantPicker") },
		ebiten.KeyC: func() { g.baize.collect() },
		ebiten.KeyU: func() { g.baize.undo() },
		ebiten.KeyV: func() {
			g.ui.Toast("Complete", fmt.Sprintf("Version %d.%d %s",
				GosoldVersionMajor,
				GosoldVersionMinor,
				GosoldVersionDate))
		},
		ebiten.KeyB: func() {
			if ebiten.IsKeyPressed(ebiten.KeyControl) {
				g.baize.loadPosition()
			} else {
				g.baize.savePosition()
			}
		},
		ebiten.KeyL:      func() { g.baize.loadPosition() },
		ebiten.KeyS:      func() { g.baize.savePosition() },
		ebiten.KeyX:      func() { ExitRequested = true },
		ebiten.KeyMenu:   func() { g.ui.ToggleNavDrawer() },
		ebiten.KeyEscape: func() { g.ui.HideActiveDrawer() },
		ebiten.KeyH: func() {
			g.settings.ShowMovableCards = !g.settings.ShowMovableCards
			if g.settings.ShowMovableCards {
				moves, _ := g.baize.darkBaize.Moves()
				if moves > 0 {
					g.ui.ToastInfo("Movable cards highlighted")
				} else {
					g.ui.ToastError("There are no movable cards")
				}
			}
		},
		ebiten.KeyM: func() {
			g.settings.AlwaysShowMovableCards = !g.settings.AlwaysShowMovableCards
			g.settings.ShowMovableCards = g.settings.AlwaysShowMovableCards
			if g.settings.AlwaysShowMovableCards {
				g.ui.ToastInfo("Movable cards always highlighted")
			}
		},
		ebiten.KeyT: func() {
			g.settings.Timer = !g.settings.Timer
		},
		ebiten.KeyA: func() {
			var AniSpeedSettings = []ui.FloatSetting{
				{Title: "Fast", Var: &g.settings.AniSpeed, Value: 0.3},
				{Title: "Normal", Var: &g.settings.AniSpeed, Value: 0.6},
				{Title: "Slow", Var: &g.settings.AniSpeed, Value: 0.9},
			}
			g.ui.ShowAniSpeedDrawer(&AniSpeedSettings)
		},
		ebiten.KeyF1: func() {
			g.baize.wikipedia()
		},
		ebiten.KeyF2: func() {
			strs := g.darker.VariantStatistics(g.baize.variant)
			strs = append(strs, " ") // n.b. can't use empty string
			strs = append(strs, "ALL VARIANTS")
			strs = append(strs, g.darker.AllStatistics()...)
			g.ui.ShowTextDrawer(strs)
		},
		ebiten.KeyF3: func() {
			var booleanSettings = []ui.BooleanSetting{
				{Title: "Power moves", Var: &g.settings.PowerMoves, Update: func() { g.baize.copySettingsToDark() }},
				{Title: "Auto collect", Var: &g.settings.AutoCollect, Update: func() { g.baize.copySettingsToDark() }},
				{Title: "Safe collect", Var: &g.settings.SafeCollect, Update: func() { g.baize.copySettingsToDark() }},
				{Title: "Show movable cards", Var: &g.settings.ShowMovableCards},
				{Title: "Colorful cards", Var: &g.settings.ColorfulCards, Update: func() { g.baize.setFlag(dirtyCardImages) }},
				{Title: "Mute sounds", Var: &g.settings.Mute, Update: func() {
					if g.settings.Mute {
						sound.SetVolume(0.0)
					} else {
						sound.SetVolume(g.settings.Volume)
					}
				}},
				{Title: "Timer", Var: &g.settings.Timer, Update: func() { g.baize.updateStatusbar() }},
				// {Title: "Mirror baize", Var: &g.settings.MirrorBaize, Update: func() {
				// 	savedUndoStack := TheGame.Baize.undoStack
				// 	TheGame.Baize.StartFreshGame()
				// 	TheGame.Baize.SetUndoStack(savedUndoStack)
				// }},
			}
			g.ui.ShowSettingsDrawer(&booleanSettings)
		},
	}

	if g.settings.LastVersionMajor != GosoldVersionMajor || g.settings.LastVersionMinor != GosoldVersionMinor {
		g.ui.Toast("Complete", fmt.Sprintf("Upgraded from %d.%d to %d.%d",
			g.settings.LastVersionMajor,
			g.settings.LastVersionMinor,
			GosoldVersionMajor,
			GosoldVersionMinor))
	}

	return g
}

func (g *Game) ExitGame() {
	// log.Println("Game.ExitGame")
	g.baize.darkBaize.Close()
	g.baize.darkBaize.Save()
	g.settings.save()
}

func (g *Game) execute(cmd any) {
	g.ui.HideActiveDrawer()
	g.ui.HideFAB()
	switch v := cmd.(type) {
	case ebiten.Key:
		if fn, ok := g.commandTable[v]; ok {
			fn()
		}
	case ui.Command:
		// a widget has sent a command
		switch v.Command {
		case "ShowVariantGroupPicker":
			g.ui.ShowVariantPickerEx(g.darker.ListVariantGroups(), "ShowVariantPicker")
		case "ShowVariantPicker":
			g.ui.ShowVariantPickerEx(g.darker.ListVariants(v.Data), "ChangeVariant")
		case "ChangeVariant":
			if v.Data == g.baize.variant {
				g.ui.ToastError(fmt.Sprintf("Already playing '%s'", v.Data))
			} else {
				g.baize.changeVariant(v.Data)
			}
		case "SaveSettings":
			g.settings.save() // save now especially if running in a browser
		default:
			log.Panic("unknown command", v.Command, v.Data)
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
	g.ui.Update()
	if ExitRequested {
		g.ExitGame()
		return errors.New("exit requested")
	}
	return nil
}

// Draw draws the current game to the given screen.
// Draw will be called based on the refresh rate of the screen (FPS).
// https://ebitencookbook.vercel.app/blog
func (g *Game) Draw(screen *ebiten.Image) {
	g.baize.draw(screen)
	g.ui.Draw(screen)
}
