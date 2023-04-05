package light

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"oddstream.games/gosold/cardid"
	"oddstream.games/gosold/dark"
	"oddstream.games/gosold/sound"
	"oddstream.games/gosold/stroke"
	"oddstream.games/gosold/ui"
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
	variant          string
	game             *Game
	darkBaize        *dark.Baize
	piles            []*pile
	cardMap          map[cardid.CardID]*card
	dirtyFlags       uint32 // what needs doing when we update
	stroke           *stroke.Stroke
	dragStart        image.Point
	dragOffset       image.Point
	windowWidth      int // the most recent window width given to Layout
	windowHeight     int // the most recent window height given to Layout
	collectRequested bool
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

func (b *baize) copySettingsToDark() {
	b.darkBaize.SetSettings(dark.BaizeSettings{
		PowerMoves:  b.game.settings.PowerMoves,
		SafeCollect: b.game.settings.SafeCollect,
	})
}

func (b *baize) eventListener(e dark.BaizeEvent, param any) {
	switch e {
	case dark.ChangedEvent:
		for _, lp := range b.piles {
			lp.copyCardsFromDark()
		}
		b.refan()
		b.updateUI()
		// juice
		// if id, ok := param.(cardid.CardID); ok {
		// 	log.Println(id.String())
		// }
		if b.game.settings.AutoCollect {
			// don't auto collect a virgin game
			if b.darkBaize.UndoStackSize() > 1 {
				if _, fmoves := b.darkBaize.Moves(); fmoves != 0 {
					b.collectRequested = true
				}
			}
		}
	case dark.LabelEvent:
		if dp, ok := param.(*dark.Pile); ok {
			for _, lp := range b.piles {
				if lp.darkPile == dp {
					lp.createPlaceholder()
				}
			}
		}
	case dark.WonEvent:
		b.game.ui.Toast("Complete", fmt.Sprintf("%s complete", b.variant))
	case dark.LostEvent:
		b.game.ui.Toast("Fail", fmt.Sprintf("Recording a lost game of %s", b.variant))
	case dark.MessageEvent:
		if str, ok := param.(string); ok {
			b.game.ui.ToastInfo(str)
		}
	}
}

func (b *baize) quiet() bool {
	if b.stroke != nil {
		return false
	}
	// for _, p := range b.piles {
	// 	for _, c := range p.cards {
	// 		if c.spinning() || c.flipping() || c.lerping() {
	// 			return false
	// 		}
	// 	}
	// }
	for _, c := range b.cardMap {
		if c.spinning() || c.flipping() || c.lerping() {
			return false
		}
	}
	return true
}

// startGame starts a new game, either an old one loaded from json,
// or a fresh game with a new seed
func (b *baize) startGame(variant string) {
	// get a new baize from dark for this variant
	var err error
	if b.darkBaize, err = b.game.darker.NewBaize(variant, b.eventListener); err != nil {
		b.game.ui.ToastError(err.Error())
		return
	}
	b.copySettingsToDark()
	b.reset()

	b.variant = variant
	b.game.ui.SetTitle(variant)
	if b.game.settings.Variant != variant {
		b.game.settings.Variant = variant
		b.game.settings.save()
	}

	// create card map
	b.cardMap = make(map[cardid.CardID]*card)
	for _, dp := range b.darkBaize.Piles() {
		for _, id := range dp.Cards() {
			b.cardMap[id.PackSuitOrdinal()] = &card{id: id}
			// card is created face up, because prone flag is not set
		}
	}
	// log.Println(len(b.cardMap), "cards in baize card map")

	// create LIGHT piles
	b.piles = []*pile{}
	for _, dp := range b.darkBaize.Piles() {
		lp := newPile(b, dp)
		b.piles = append(b.piles, lp)
		lp.copyCardsFromDark()
		lp.createPlaceholder()
	}
	// log.Println(len(b.piles), "piles created")

	if b.game.settings.MirrorBaize {
		b.mirrorSlots()
	}

	sound.Play("Fan")
	b.dirtyFlags = dirtyAll
	b.updateUI()
}

func (b *baize) newDeal() {
	b.stopSpinning()
	if ok, err := b.darkBaize.NewDeal(); !ok {
		b.game.ui.ToastError(err.Error())
	} else {
		sound.Play("Fan")
	}
}

func (b *baize) restartDeal() {
	if ok, err := b.darkBaize.RestartDeal(); !ok {
		b.game.ui.ToastError(err.Error())
	} else {
		sound.Play("Fan")
	}
}

func (b *baize) changeVariant(variant string) {
	b.darkBaize.Save()
	b.startGame(variant)
	b.darkBaize.Load()
}

func (b *baize) collect() {
	if _, fmoves := b.darkBaize.Moves(); fmoves != 0 {
		// even with fmoves > 0, SafeCollect may mean no cards are collected
		if b.darkBaize.Collect() > 0 {
			sound.Play("Shove")
		}
	}
}

func (b *baize) undo() {
	// temporarily disable autocollect
	var saved bool = b.game.settings.AutoCollect
	b.game.settings.AutoCollect = false
	if ok, err := b.darkBaize.Undo(); !ok {
		b.game.ui.ToastError(err.Error())
	} else {
		sound.Play("Fan")
	}
	b.game.settings.AutoCollect = saved
}

func (b *baize) loadPosition() {
	if ok, err := b.darkBaize.LoadPosition(); !ok {
		b.game.ui.ToastError(err.Error())
		return
	}
}

func (b *baize) savePosition() {
	if ok, err := b.darkBaize.SavePosition(); !ok {
		b.game.ui.ToastError(err.Error())
		return
	}
	b.game.ui.ToastInfo("Position bookmarked")
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
// func (b *baize) findHighestCardAt(pt image.Point) *card {
// 	for _, p := range b.piles {
// 		for _, c := range p.cards {
// 			if pt.In(c.screenRect()) {
// 				return c
// 			}
// 		}
// 	}
// 	return nil
// }

func (b *baize) largestIntersection(c *card) *pile {
	var largestArea int = 0
	var largest *pile = nil
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
			largest = p
		}
	}
	return largest
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
	for _, c := range b.cardMap {
		c.startSpinning()
	}
	// for _, p := range b.piles {
	// 	// use a method expression, which yields a function value with a regular first parameter taking the place of the receiver
	// 	p.applyToCards((*card).startSpinning)
	// }
}

// stopSpinning tells all the cards to stop spinning and return to their upright position
func (b *baize) stopSpinning() {
	for _, c := range b.cardMap {
		c.stopSpinning()
	}
	// for _, p := range b.piles {
	// 	// use a method expression, which yields a function value with a regular first parameter taking the place of the receiver
	// 	p.applyToCards((*card).stopSpinning)
	// }
	b.setFlag(dirtyCardPositions)
}

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

	TopMargin = ui.ToolbarHeight + CardHeight/3
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
		}
		if b.flagSet(dirtyCardImages) {
			b.game.createCardImages()
		}
		if b.flagSet(dirtyPilePositions) {
			for _, p := range b.piles {
				p.setBaizePos(image.Point{
					X: LeftMargin + (p.slot.X * (CardWidth + PilePaddingX)),
					Y: TopMargin + (p.slot.Y * (CardHeight + PilePaddingY)),
				})
			}
		}
		if b.flagSet(dirtyPileBackgrounds) {
			if !(CardWidth == 0 || CardHeight == 0) {
				for i := range b.piles {
					b.piles[i].createPlaceholder()
				}
			}
		}
		if b.flagSet(dirtyWindowSize) {
			b.game.ui.Layout(outsideWidth, outsideHeight)
		}
		if b.flagSet(dirtyCardPositions) {
			for _, p := range b.piles {
				p.scrunch()
			}
		}
		b.dirtyFlags = 0
	}

	return outsideWidth, outsideHeight
}

// foreachCard applys a function to each card
// func (b *baize) foreachCard(fn func(*card)) {
// 	for _, p := range b.piles {
// 		for _, c := range p.cards {
// 			fn(c)
// 		}
// 	}
// }

// ApplyToTail applies a method func to this card and all the others after it in the tail
func (b *baize) applyToTail(tail []*card, fn func(*card)) {
	// https://golang.org/ref/spec#Method_expressions
	// (*Card).CancelDrag yields a function with the signature func(*Card)
	// fn passed as a method expression so add the receiver explicitly
	for _, c := range tail {
		fn(c)
	}
}

// DragTailBy repositions all the cards in the tail: dx, dy is the position difference from the start of the drag
func (b *baize) dragTailBy(tail []*card, dx, dy int) {
	// println("Baize.DragTailBy(", dx, dy, ")")
	for _, c := range tail {
		c.dragBy(dx, dy)
	}
}

func (b *baize) startTailDrag(tail []*card) {
	// hiding the mouse cursor creates flickering when tapping
	// ebiten.SetCursorMode(ebiten.CursorModeHidden)
	b.applyToTail(tail, (*card).startDrag)
}

func (b *baize) stopTailDrag(tail []*card) {
	// ebiten.SetCursorMode(ebiten.CursorModeVisible)
	b.applyToTail(tail, (*card).stopDrag)
}

func (b *baize) cancelTailDrag(tail []*card) {
	// ebiten.SetCursorMode(ebiten.CursorModeVisible)
	b.applyToTail(tail, (*card).cancelDrag)
}

func (b *baize) findCard(cid cardid.CardID) *card {
	for _, c := range b.cardMap {
		if c.id == cid {
			return c
		}
	}
	// for _, p := range b.piles {
	// 	for _, c := range p.cards {
	// 		if c.id == cid {
	// 			return c
	// 		}
	// 	}
	// }
	return nil
}

func (b *baize) strokeStart(v stroke.StrokeEvent) {
	b.stroke = v.Stroke

	if con := b.game.ui.FindContainerAt(v.X, v.Y); con != nil {
		if w := con.FindWidgetAt(v.X, v.Y); w != nil {
			b.stroke.SetDraggedObject(w)
		} else {
			con.StartDrag()
			b.stroke.SetDraggedObject(con)
		}
	} else {
		pt := image.Pt(v.X, v.Y)
		if c := b.findLowestCardAt(pt); c != nil {
			if c.lerping() {
				// TheGame.UI.Toast("Glass", "Confusing to move a moving card")
				v.Stroke.Cancel()
			} else {
				tail := c.pile.makeTail(c)
				b.startTailDrag(tail)
				b.stroke.SetDraggedObject(tail) // TODO use card.id instead
			}
		} else {
			if p := b.findPileAt(pt); p != nil {
				b.stroke.SetDraggedObject(p)
			} else {
				if b.startDrag() {
					b.stroke.SetDraggedObject(b)
				} else {
					v.Stroke.Cancel()
				}
			}
		}
	}
}

func (b *baize) strokeMove(v stroke.StrokeEvent) {
	if v.Stroke.DraggedObject() == nil {
		return
		// log.Panic("*** move stroke with nil dragged object ***")
	}
	// for _, p := range b.piles {
	// 	p.target = false
	// }
	switch obj := v.Stroke.DraggedObject().(type) {
	case ui.Containery:
		obj.DragBy(v.Stroke.PositionDiff())
	case ui.Widgety:
		obj.Parent().DragBy(v.Stroke.PositionDiff())
	case cardid.CardID:
		c := b.findCard(obj)
		tail := c.pile.makeTail(c)
		dx, dy := v.Stroke.PositionDiff()
		b.dragTailBy(tail, dx, dy)
	case []*card:
		dx, dy := v.Stroke.PositionDiff()
		b.dragTailBy(obj, dx, dy)
		// if c, ok := v.Stroke.DraggedObject().(*Card); ok {
		// 	if p := b.LargestIntersection(c); p != nil {
		// 		p.target = true
		// 	}
		// }
	case *pile:
		// do nothing
	case *baize:
		b.dragBy(v.Stroke.PositionDiff())
	default:
		log.Panic("*** unknown move dragging object ***")
	}
}

func (b *baize) strokeStop(v stroke.StrokeEvent) {
	if v.Stroke.DraggedObject() == nil {
		return
		// log.Panic("*** stop stroke with nil dragged object ***")
	}
	switch obj := v.Stroke.DraggedObject().(type) {
	case ui.Containery:
		obj.StopDrag()
	case ui.Widgety:
		obj.Parent().StopDrag()
	// case cardid.CardID:
	// 	c := b.findCard(obj)
	// 	tail := c.pile.makeTail(c)
	case []*card:
		tail := obj  // alias for readability
		c := tail[0] // for readability
		if c.wasDragged() {
			if dst := b.largestIntersection(c); dst == nil {
				// println("no intersection for", c.String())
				b.cancelTailDrag(tail)
			} else {
				src := c.pile
				if ok, err := b.darkBaize.CardDragged(src.darkPile, c.id, dst.darkPile); !ok {
					b.game.ui.ToastError(err.Error())
					b.cancelTailDrag(tail)
				} else {
					sound.Play("Slide")
					b.stopTailDrag(tail)
				}
			}
		}
	case *pile:
		// do nothing
	case *baize:
		// println("stop dragging baize")
		b.stopDrag()
	default:
		log.Panic("*** stop dragging unknown object ***")
	}
}

func (b *baize) strokeCancel(v stroke.StrokeEvent) {
	if v.Stroke.DraggedObject() == nil {
		log.Print("*** cancel stroke with nil dragged object ***")
		return
	}
	switch obj := v.Stroke.DraggedObject().(type) { // type switch
	case ui.Containery:
		obj.CancelDrag()
	case ui.Widgety:
		obj.Parent().CancelDrag()
	case []*card:
		b.cancelTailDrag(obj)
	case *pile:
		// p := v.Stroke.DraggedObject().(*Pile)
		// println("stop dragging pile", p.Class)
		// do nothing
	case *baize:
		// println("stop dragging baize")
		b.stopDrag()
	default:
		log.Panic("*** cancel dragging unknown object ***")
	}
}

func (b *baize) strokeTap(v stroke.StrokeEvent) {
	// stroke sends a tap event, and later sends a cancel event
	// println("Baize.NotifyCallback() tap", v.X, v.Y)
	switch obj := v.Stroke.DraggedObject().(type) {
	case ui.Containery:
		obj.Tapped()
	case ui.Widgety:
		obj.Tapped()
	case []*card:
		// offer TailTapped to the script first
		// to implement things like Stock.TailTapped
		// if the script doesn't want to do anything, it can call pile.vtable.TailTapped
		// which will either ignore it (eg Foundation, Discard)
		// or use Pile.DefaultTailTapped

		// obj is a tail of light cards, but we need a tail of dark cards
		// also, if have run the solver, then the card.darkCard objects will no longer exist

		if b.darkBaize.CardTapped(obj[0].id) {
			sound.Play("Slide")
		} else {
			sound.Play("Glass")
		}
	case *pile:
		if b.darkBaize.PileTapped(obj.darkPile) {
			sound.Play("Shove")
		} else {
			sound.Play("Glass")
		}
	case *baize:
		pt := image.Pt(v.X, v.Y)
		// 	// a tap outside any open ui drawer (ie on the baize) closes the drawer
		if con := b.game.ui.VisibleDrawer(); con != nil && !pt.In(image.Rect(con.Rect())) {
			con.Hide()
		}
	default:
		log.Panic("*** tap unknown object ***")
	}
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (b *baize) NotifyCallback(v stroke.StrokeEvent) {
	switch v.Event {
	case stroke.Start:
		b.strokeStart(v)
	case stroke.Move:
		b.strokeMove(v)
	case stroke.Stop:
		b.strokeStop(v)
	case stroke.Cancel:
		b.strokeCancel(v)
	case stroke.Tap:
		b.strokeTap(v)
	default:
		log.Panic("*** unknown stroke event ***", v.Event)
	}
}

func (b *baize) updateToolbar() {
	undos := b.darkBaize.UndoStackSize()
	b.game.ui.EnableWidget("toolbarUndo", undos > 1)
	_, fmoves := b.darkBaize.Moves()
	b.game.ui.EnableWidget("toolbarCollect", fmoves > 0)
}

func (b *baize) updateStatusbar() {
	b.game.ui.SetStock(b.darkBaize.StockLen())
	b.game.ui.SetWaste(b.darkBaize.WasteLen())
	b.game.ui.SetMiddle(fmt.Sprintf("MOVES: %d", b.darkBaize.UndoStackSize()-1))
	percent, fpercent := b.darkBaize.PercentComplete()
	if DebugMode {
		b.game.ui.SetPercent2(fmt.Sprintf("%d%%/%d%%", percent, fpercent))
	} else {
		b.game.ui.SetPercent(percent)
	}
}

func (b *baize) updateDrawers() {
	b.game.ui.EnableWidget("restartDeal", b.darkBaize.UndoStackSize() > 1)
	b.game.ui.EnableWidget("gotoBookmark", b.darkBaize.Bookmark() > 0)
}

func (b *baize) updateFAB() {
	b.game.ui.HideFAB()
	if b.darkBaize.Complete() {
		b.game.ui.AddButtonToFAB("star", ebiten.KeyN)
		b.startSpinning()
	} else if b.darkBaize.Conformant() {
		b.game.ui.AddButtonToFAB("done_all", ebiten.KeyC)
	} else if moves, _ := b.darkBaize.Moves(); moves == 0 {
		b.game.ui.ToastError("No movable cards")
		b.game.ui.AddButtonToFAB("star", ebiten.KeyN)
		b.game.ui.AddButtonToFAB("restore", ebiten.KeyR)
		if b.darkBaize.Bookmark() > 0 {
			b.game.ui.AddButtonToFAB("bookmark", ebiten.KeyL)
		}
	}
}

func (b *baize) updateUI() {
	b.updateToolbar()
	b.updateDrawers()
	b.updateStatusbar()
	b.updateFAB()
}

func (b *baize) update() error {

	if b.stroke == nil {
		stroke.StartStroke(b) // this will set b.stroke when "start" received
	} else {
		b.stroke.Update()
		if b.stroke.IsReleased() || b.stroke.IsCancelled() {
			b.stroke = nil
		}
	}

	for _, p := range b.piles {
		p.update()
	}

	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if inpututil.IsKeyJustReleased(k) {
			b.game.execute(k)
		}
	}

	if b.collectRequested && b.quiet() {
		b.collect()
		b.collectRequested = false
	}

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
