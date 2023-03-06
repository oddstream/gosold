package light

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosold/cardid"
	"oddstream.games/gosold/dark"
	"oddstream.games/gosold/sound"
	"oddstream.games/gosold/stroke"
)

const (
	dirtyWindowSize = 1 << iota
	dirtyPilePositions
	dirtyCardSizes
	dirtyCardImages
	dirtyPileBackgrounds
	dirtyCardPositions
)

const dirtyAll uint32 = 0xFFFF

type baize struct {
	variant      string
	game         *Game
	darker       dark.Darker
	darkBaize    *dark.Baize
	piles        []*pile
	cardMap      map[cardid.CardID]*card
	dirtyFlags   uint32 // what needs doing when we Update
	stroke       *stroke.Stroke
	dragStart    image.Point
	dragOffset   image.Point
	WindowWidth  int // the most recent window width given to Layout
	WindowHeight int // the most recent window height given to Layout
}

func newBaize(g *Game) *baize {
	return &baize{game: g, dirtyFlags: dirtyAll}
}

func (b *baize) flagSet(flag uint32) bool {
	return b.dirtyFlags&flag == flag
}

func (b *baize) setFlag(flag uint32) {
	b.dirtyFlags |= flag
}

func (b *baize) refan() {
	b.setFlag(dirtyCardPositions)
}

func (b *baize) updateCardsAndLabels() {
	for i, dp := range b.darkBaize.Piles() {
		lp := b.piles[i]
		lp.updateCardsAndLabel(dp)
	}
}

func (b *baize) startFreshGame(variant string) {
	// get a new baize from dark for this variant
	var err error
	if b.darkBaize, err = b.darker.NewBaize(variant); err != nil {
		// toast
		return
	}
	// copy piles from dark to light
	// the pile structure will remain constant
	// the labels and cards in each pile will change
	// so are rebuilt after every move
	b.piles = []*pile{}
	b.cardMap = make(map[cardid.CardID]*card)
	for _, dp := range b.darkBaize.Piles() {
		b.piles = append(b.piles, newPile(b, dp))
		for _, dc := range dp.Cards() {
			c := newCard(dc)
			b.cardMap[c.darkCard.ID().PackSuitOrdinal()] = c
		}
	}

	b.updateCardsAndLabels()
	b.variant = variant
	sound.Play("Fan")
	b.dirtyFlags = dirtyAll
}

func (b *baize) changeVariant(variant string) {
}

func (b *baize) maxSlotX() int {
	var maxX int
	for _, p := range b.darkBaize.Piles() {
		if p.Slot().X > maxX {
			maxX = p.Slot().X
		}
	}
	return maxX
}

// ScaleCards calculates new width/height of cards and margins
// returns true if changes were made
func (b *baize) ScaleCards() bool {

	// const (
	// 	DefaultRatio = 1.444
	// 	BridgeRatio  = 1.561
	// 	PokerRatio   = 1.39
	// 	OpsoleRatio  = 1.5556 // 3.5/2.25
	// )

	var OldWidth = CardWidth
	var OldHeight = CardHeight

	var maxX int = b.maxSlotX()

	var slotWidth, slotHeight float64
	slotWidth = float64(b.WindowWidth) / float64(maxX+2)
	slotHeight = slotWidth * b.game.settings.CardRatio

	PilePaddingX = int(slotWidth / 10)
	CardWidth = int(slotWidth) - PilePaddingX
	PilePaddingY = int(slotHeight / 10)
	CardHeight = int(slotHeight) - PilePaddingY

	TopMargin = /* ui.ToolbarHeight */ 48 + CardHeight/3
	LeftMargin = (CardWidth / 2) + PilePaddingX

	CardCornerRadius = float64(CardWidth) / 10.0 // same as lsol
	return CardWidth != OldWidth || CardHeight != OldHeight
}

func (b *baize) layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth == 0 || outsideHeight == 0 {
		log.Println("Baize.Layout called with zero dimension")
		return outsideWidth, outsideHeight
	}

	if outsideWidth != b.WindowWidth {
		b.setFlag(dirtyWindowSize | dirtyCardSizes | dirtyPileBackgrounds | dirtyPilePositions | dirtyCardPositions)
		b.WindowWidth = outsideWidth
	}
	if outsideHeight != b.WindowHeight {
		b.setFlag(dirtyWindowSize | dirtyCardPositions)
		b.WindowHeight = outsideHeight
	}

	if b.dirtyFlags != 0 {
		if b.flagSet(dirtyCardSizes) {
			if b.ScaleCards() {
				b.setFlag(dirtyCardImages | dirtyPilePositions | dirtyPileBackgrounds)
			}
			// b.clearFlag(dirtyCardSizes)
		}
		if b.flagSet(dirtyCardImages) {
			b.game.createCardImages()
			// b.clearFlag(dirtyCardImages)
		}
		if b.flagSet(dirtyPilePositions) {
			for _, p := range b.piles {
				p.setBaizePos(image.Point{
					X: LeftMargin + (p.darkPile.Slot().X * (CardWidth + PilePaddingX)),
					Y: TopMargin + (p.darkPile.Slot().Y * (CardHeight + PilePaddingY)),
				})
			}
			// b.clearFlag(dirtyPilePositions)
		}
		if b.flagSet(dirtyPileBackgrounds) {
			if !(CardWidth == 0 || CardHeight == 0) {
				for i, p := range b.darkBaize.Piles() {
					if !p.Hidden() {
						b.piles[i].createPlaceholder()
					}
				}
			}
			// b.clearFlag(dirtyPileBackgrounds)
		}
		if b.flagSet(dirtyWindowSize) {
			// b.game.ui.Layout(outsideWidth, outsideHeight)
			// b.clearFlag(dirtyWindowSize)
		}
		if b.flagSet(dirtyCardPositions) {
			for _, p := range b.piles {
				p.Scrunch()
			}
			// b.clearFlag(dirtyCardPositions)
		}
		b.dirtyFlags = 0
	}

	return outsideWidth, outsideHeight
}

func (b *baize) update() error {
	// TODO stroke
	for _, p := range b.piles {
		p.update()
	}
	// TODO keys
	return nil
}

func (b *baize) draw(screen *ebiten.Image) {
	screen.Fill(ExtendedColors["BaizeGreen"])
	for _, p := range b.piles {
		p.draw(screen)
	}
	// TODO draw static, lerping, dragging cards
}
