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
	darkBaize    *dark.Baize
	piles        []*pile
	cardMap      map[cardid.CardID]*card
	dirtyFlags   uint32 // what needs doing when we Update
	stroke       *stroke.Stroke
	dragStart    image.Point
	dragOffset   image.Point
	windowWidth  int // the most recent window width given to Layout
	windowHeight int // the most recent window height given to Layout
}

func newBaize(g *Game) *baize {
	return &baize{game: g, dirtyFlags: dirtyAll}
}

func (b *baize) reset() {
	b.stopSpinning()
}

func (b *baize) flagSet(flag uint32) bool {
	return b.dirtyFlags&flag == flag
}

func (b *baize) setFlag(flag uint32) {
	b.dirtyFlags |= flag
}

// startGame starts a new game, either an old one loaded from json,
// or a fresh game with a new seed
func (b *baize) startGame(variant string) {
	// get a new baize from dark for this variant
	var err error
	if b.darkBaize, err = b.game.darker.NewBaize(variant); err != nil {
		// TODO toast
		println(err)
		return
	}

	b.reset()

	b.variant = variant

	// create card map
	b.cardMap = make(map[cardid.CardID]*card)
	for _, dp := range b.darkBaize.Piles() {
		for _, dc := range dp.Cards() {
			b.cardMap[dc.ID().PackSuitOrdinal()] = newCard(dc)
		}
	}
	println(len(b.cardMap), "cards in baize card map")

	// create LIGHT piles
	b.piles = []*pile{}
	for _, dp := range b.darkBaize.Piles() {
		lp := newPile(b, dp)
		b.piles = append(b.piles, lp)
		lp.updateCards()
		lp.createPlaceholder()
	}
	println(len(b.piles), "piles created")

	if b.game.settings.MirrorBaize {
		b.mirrorSlots()
	}

	sound.Play("Fan")
	b.dirtyFlags = dirtyAll
}

func (b *baize) newDeal() {
	b.darkBaize.NewDeal()
	for _, lp := range b.piles {
		lp.updateCards()
		lp.createPlaceholder()
	}
}

func (b *baize) restartDeal() {
	if ok, err := b.darkBaize.RestartDeal(); !ok {
		// TODO toast err
		println(err)
		return
	}
	for _, lp := range b.piles {
		lp.updateCards()
		lp.createPlaceholder()
	}
}

func (b *baize) changeVariant(variant string) {
	b.darkBaize.Save()
	b.startGame(variant)
}

func (b *baize) collect() {
	b.darkBaize.Collect(b.game.settings.SafeCollect)
	for _, lp := range b.piles {
		lp.updateCards()
		lp.createPlaceholder()
	}
}

func (b *baize) undo() {
	if ok, err := b.darkBaize.Undo(); !ok {
		// TODO toast err
		println(err)
		return
	}
	for _, lp := range b.piles {
		lp.updateCards()
		lp.createPlaceholder()
	}
}

func (b *baize) loadPosition() {
	if ok, err := b.darkBaize.LoadPosition(); !ok {
		// TODO toast err
		println(err)
		return
	}
	for _, lp := range b.piles {
		lp.updateCards()
		lp.createPlaceholder()
	}
}

func (b *baize) savePosition() {
	if ok, err := b.darkBaize.SavePosition(); !ok {
		// TODO toast err
		println(err)
		return
	}
	// TODO recycles may have changed, so may need to recreate Stock placeholder
}

// findPileAt finds the Pile under the mouse position
func (b *baize) findPileAt(pt image.Point) *pile {
	for _, p := range b.piles {
		if pt.In(p.screenRect()) {
			return p
		}
	}
	return nil
}

// findLowestCardAt finds the bottom-most Card under the mouse position
func (b *baize) findLowestCardAt(pt image.Point) *card {
	for _, p := range b.piles {
		for i := len(p.cards) - 1; i >= 0; i-- {
			c := p.cards[i]
			if pt.In(c.screenRect()) {
				return c
			}
		}
	}
	return nil
}

// findHighestCardAt finds the top-most Card under the mouse position
func (b *baize) findHighestCardAt(pt image.Point) *card {
	for _, p := range b.piles {
		for _, c := range p.cards {
			if pt.In(c.screenRect()) {
				return c
			}
		}
	}
	return nil
}

func (b *baize) largestIntersection(c *card) *pile {
	var largestArea int = 0
	var pile *pile = nil
	cardRect := c.baizeRect()
	for _, p := range b.piles {
		if p == c.pile {
			continue
		}
		pileRect := p.fannedBaizeRect()
		intersectRect := pileRect.Intersect(cardRect)
		area := intersectRect.Dx() * intersectRect.Dy()
		if area > largestArea {
			largestArea = area
			pile = p
		}
	}
	return pile
}

// startDrag return true if the Baize can be dragged
func (b *baize) startDrag() bool {
	b.dragStart = b.dragOffset
	return true
}

// dragBy move ('scroll') the Baize by dragging it
// dx, dy is the difference between where the drag started and where the cursor is now
func (b *baize) dragBy(dx, dy int) {
	b.dragOffset.X = b.dragStart.X + dx
	if b.dragOffset.X > 0 {
		b.dragOffset.X = 0 // DragOffsetX should only ever be 0 or -ve
	}
	b.dragOffset.Y = b.dragStart.Y + dy
	if b.dragOffset.Y > 0 {
		b.dragOffset.Y = 0 // DragOffsetY should only ever be 0 or -ve
	}
}

// stopDrag stop dragging the Baize
func (b *baize) stopDrag() {
	b.setFlag(dirtyCardPositions)
}

// startSpinning tells all the cards to start spinning
func (b *baize) startSpinning() {
	for _, p := range b.piles {
		// use a method expression, which yields a function value with a regular first parameter taking the place of the receiver
		p.applyToCards((*card).startSpinning)
	}
}

// stopSpinning tells all the cards to stop spinning and return to their upright position
func (b *baize) stopSpinning() {
	for _, p := range b.piles {
		// use a method expression, which yields a function value with a regular first parameter taking the place of the receiver
		p.applyToCards((*card).stopSpinning)
	}
	b.setFlag(dirtyCardPositions)
}

// TODO input

func (b *baize) mirrorSlots() {
	/*
		0 1 2 3 4 5
		5 4 3 2 1 0

		0 1 2 3 4
		4 3 2 1 0
	*/
	var minX int = 32767
	var maxX int = 0
	for _, p := range b.piles {
		if p.slot.X < 0 {
			continue // ignore hidden pile
		}
		if p.slot.X < minX {
			minX = p.slot.X
		}
		if p.slot.X > maxX {
			maxX = p.slot.X
		}
	}
	for _, p := range b.piles {
		slot := p.slot
		if slot.X < 0 {
			continue // ignore hidden pile
		}
		p.slot = image.Point{X: maxX - slot.X + minX, Y: slot.Y}
		switch p.fanType {
		case dark.FAN_RIGHT:
			p.fanType = dark.FAN_LEFT
		case dark.FAN_LEFT:
			p.fanType = dark.FAN_RIGHT
		case dark.FAN_RIGHT3:
			p.fanType = dark.FAN_LEFT3
		case dark.FAN_LEFT3:
			p.fanType = dark.FAN_RIGHT3
		}
	}
}

func (b *baize) refan() {
	b.setFlag(dirtyCardPositions)
}

func (b *baize) maxSlotX() int {
	// nb use local copy of slot, not darkPile.Slot()
	var maxX int
	for _, p := range b.piles {
		if p.slot.X > maxX {
			maxX = p.slot.X
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
	slotWidth = float64(b.windowWidth) / float64(maxX+2)
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

	if outsideWidth != b.windowWidth {
		b.setFlag(dirtyWindowSize | dirtyCardSizes | dirtyPileBackgrounds | dirtyPilePositions | dirtyCardPositions)
		b.windowWidth = outsideWidth
	}
	if outsideHeight != b.windowHeight {
		b.setFlag(dirtyWindowSize | dirtyCardPositions)
		b.windowHeight = outsideHeight
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
					X: LeftMargin + (p.slot.X * (CardWidth + PilePaddingX)),
					Y: TopMargin + (p.slot.Y * (CardHeight + PilePaddingY)),
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
				p.scrunch()
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
	for _, p := range b.piles {
		p.drawStaticCards(screen)
	}
	for _, p := range b.piles {
		p.drawAnimatingCards(screen)
	}
	for _, p := range b.piles {
		p.drawDraggingCards(screen)
	}
}
