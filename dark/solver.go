package dark

import "log"

func (b *Baize) createChild() *Baize {
	if b.crc == 0 {
		println("parent crc not set")
		b.crc = b.calcCRC()
	}
	child := &Baize{
		dark:       b.dark,
		variant:    b.variant,
		script:     b.script,
		cardCount:  b.cardCount,
		piles:      make([]*Pile, 0, len(b.piles)),
		recycles:   b.recycles,
		percent:    b.percent,
		crc:        b.crc,
		depth:      b.depth + 1,
		parent:     b,
		children:   []*Baize{},
		tappedCard: 0,
		bookmark:   0,
		moves:      0,
		fmoves:     0,
		undoStack:  []*savableBaize{},
		fnNotify: func(evt BaizeEvent, data any) {
		},
		BaizeSettings: b.BaizeSettings,
	}

	// create piles and set script aliases for foundations, waste &c
	child.script.SetBaize(child)
	child.script.BuildPiles()

	// clone cards in each pile
	for i, src := range b.piles {
		dst := child.piles[i]
		// clone cards in each pile
		for _, c := range src.cards {
			dst.cards = append(dst.cards, &Card{id: c.id, pile: dst})
		}
		if len(src.cards) != len(dst.cards) {
			println("createChild len cards mismatch")
		}
	}
	if len(child.piles) != len(b.piles) {
		println("createChild len piles mismatch")
	}
	for i := range b.piles {
		for j := range b.piles[i].cards {
			if b.piles[i].cards[j].id != child.piles[i].cards[j].id {
				println("card mismatch @ pile", i, "card", j)
			}
		}
	}
	if child.calcCRC() != b.calcCRC() {
		println("calcCRC mismatch child:", child.calcCRC(), "parent:", b.calcCRC())
	}
	child.crc = child.calcCRC()
	if child.crc != b.crc {
		println("createChild crc mismatch child:", child.crc, "parent:", b.crc)
	}
	child.findTapTargets()
	// child.tappedCard will be added later
	return child
}

func (b *Baize) findCRC(crc uint32) bool {
	// if b.crc == crc {
	// 	return true
	// }
	// for _, child := range b.children {
	// 	if child.findCRC(crc) {
	// 		return true
	// 	}
	// }
	return false
}

func countTargets(b *Baize) int {
	targets := 0
	for _, p := range b.piles {
		for _, c := range p.cards {
			targets += len(c.tapTargets)
			// for i, tt := range c.tapTargets {
			// 	println("[", c.String(), "]", i, tt.dst.category, tt.weight)
			// }
		}
	}
	return targets
}

func solve(root *Baize, b *Baize, depth int) {
	if b.depth > depth || b.percent >= 100 {
		println("limit reached", b.depth, b.percent)
		return
	}
	if b.moves == 0 {
		println("stuck")
		return
	}
	println("depth:", b.depth, "moves:", b.moves)

	println("parent has", countTargets(b), "tap targets")

	root.crc = root.calcCRC()

	b1 := b.createChild()
	println("child has", countTargets(b1), "tap targets")

	// for each card with tap targets
	for _, p := range b1.piles {
		for _, c := range p.cards {
			for nTapTarget := range c.tapTargets {
				tail := p.makeTail(c)
				oldCRC := b1.calcCRC()
				// move tail to a pre-determined destination
				b1.script.TailTapped(tail, nTapTarget)
				if b1.calcCRC() == oldCRC {
					println("nothing changed!")
				} else {
					// println("baize changed")
					// the tap caused the baize to change
					b1.script.AfterMove()
					b1.percent = b1.percentComplete()
					b1.tappedCard = c.id
					b.children = append(b.children, b1)
					// solve(root, b1, depth)
				}
			}
		}
	}

}

func display(b *Baize) {
	log.Printf("depth %d complete %d%%", b.depth, b.percent)
	for _, b1 := range b.children {
		display(b1)
	}
}

func (b *Baize) Solve(depth int) {
	b.children = nil
	solve(b, b, depth)
	display(b)
}
