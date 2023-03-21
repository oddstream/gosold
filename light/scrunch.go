package light

import (
	"oddstream.games/gosold/dark"
)

// func (b *Baize) FindBuddyPiles() {
// 	for _, p1 := range b.piles {
// 		switch (p1).(type) {
// 		case *Tableau, *Reserve:
// 			p1.buddyPos = image.Point{0, 0}
// 			for _, p2 := range b.piles {
// 				switch (p2).(type) {
// 				case *Tableau, *Reserve:
// 					switch p1.fanType {
// 					case FAN_DOWN:
// 						if p1.slot.X == p2.slot.X && p2.slot.Y > p1.slot.Y {
// 							p1.buddyPos = p2.pos
// 						}
// 					case FAN_LEFT:
// 						if p1.slot.Y == p2.slot.Y && p2.slot.X < p1.slot.X {
// 							p1.buddyPos = p2.pos
// 						}
// 					case FAN_RIGHT:
// 						if p1.slot.Y == p2.slot.Y && p2.slot.X > p1.slot.X {
// 							p1.buddyPos = p2.pos
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// func (b *Baize) CalcScrunchDims(w, h int) {
// 	for _, pile := range b.piles {
// 		switch (pile).(type) {
// 		case *Tableau, *Reserve:
// 			switch pile.FanType() {
// 			case FAN_DOWN:
// 				if pile.buddyPos.Y != 0 {
// 					pile.scrunchDims.Y = pile.buddyPos.Y - pile.pos.Y
// 				} else {
// 					// baize->dragOffset is always -ve
// 					pile.scrunchDims.Y = h - pile.pos.Y + util.Abs(b.dragOffset.Y)
// 				}
// 			case FAN_LEFT:
// 				if pile.buddyPos.X != 0 {
// 					pile.scrunchDims.X = pile.buddyPos.X - pile.pos.X
// 				} else {
// 					pile.scrunchDims.X = pile.pos.X
// 				}
// 			case FAN_RIGHT:
// 				if pile.buddyPos.X != 0 {
// 					pile.scrunchDims.X = pile.buddyPos.X - pile.pos.X
// 				} else {
// 					// baize->dragOffset is always -ve
// 					pile.scrunchDims.X = w - pile.pos.X + util.Abs(b.dragOffset.X)
// 				}
// 			}
// 			pile.fanFactor = DefaultFanFactor[pile.fanType]
// 		}
// 	}
// }

// sizeWithFanFactor calculates the width or height this pile would be if it had a specified fan factor
func (p *pile) sizeWithFanFactor(fanFactor float64) int {
	var max int
	switch p.fanType {
	case dark.FAN_DOWN:
		for i := 0; i < len(p.cards)-1; i++ {
			c := p.cards[i]
			if c.pile.baize.darkBaize.IsCardProne(c.id) {
				max += int(float64(CardHeight) / CARD_BACK_FAN_FACTOR)
			} else {
				max += int(float64(CardHeight) / fanFactor)
			}
		}
		max += CardHeight
	case dark.FAN_LEFT, dark.FAN_RIGHT:
		for i := 0; i < len(p.cards)-1; i++ {
			c := p.cards[i]
			if c.pile.baize.darkBaize.IsCardProne(c.id) {
				max += int(float64(CardWidth) / CARD_BACK_FAN_FACTOR)
			} else {
				max += int(float64(CardWidth) / fanFactor)
			}
		}
		max += CardWidth
	}
	return max
}

// scrunch prepares to refan cards after Push() or Pop(), adjusting the amount of overlap to try to keep them fitting on the screen
// only Scrunch piles with fanType LEFT/RIGHT/UP/DOWN, ignore the waste-style piles and those that do not fan
func (p *pile) scrunch() {

	p.fanFactor = defaultFanFactor[p.fanType]

	if len(p.cards) < 2 {
		p.refan()
		return
	}

	var maxPileSize int
	switch p.fanType {
	case dark.FAN_DOWN:
		// baize->dragOffset is always -ve
		// statusbar height is 24
		// maxPileSize = TheGame.Baize.WindowHeight - scpos.Y + util.Abs(TheGame.Baize.dragOffset.Y)
		maxPileSize = p.baize.windowHeight - p.screenPos().Y + (CardHeight / 2)
	case dark.FAN_LEFT:
		maxPileSize = p.screenPos().X
	case dark.FAN_RIGHT:
		// baize->dragOffset is always -ve
		// maxPileSize = TheGame.Baize.WindowWidth - scpos.X + util.Abs(TheGame.Baize.dragOffset.X)
		maxPileSize = p.baize.windowWidth - p.screenPos().X
	}
	if maxPileSize == 0 {
		// this pile doesn't need scrunching
		p.refan()
		return
	}

	var nloops int
	var fanFactor float64
	for fanFactor = defaultFanFactor[p.fanType]; fanFactor < 7.0; fanFactor += 0.1 {
		size := p.sizeWithFanFactor(fanFactor)
		switch p.fanType {
		case dark.FAN_DOWN:
			if size < maxPileSize {
				goto exitloop
			}
		case dark.FAN_LEFT, dark.FAN_RIGHT:
			if size < maxPileSize {
				goto exitloop
			}
		default:
			goto exitloop
		}
		nloops++
	}
exitloop:
	p.fanFactor = fanFactor
	p.refan()
}
