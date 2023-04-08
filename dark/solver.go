package dark

import (
	"log"

	"oddstream.games/gosold/cardid"
)

/*
func (b *Baize) createChild() *Baize {
	if b.crc == 0 {
		println("parent crc not set")
		b.crc = b.calcCRC()
	}
	child := &Baize{
		dark:       b.dark,
		variant:    b.variant,
		script:     variants[b.variant],
		cardCount:  b.cardCount,
		piles:      make([]*Pile, 0, len(b.piles)),
		recycles:   b.recycles,
		percent:    b.percent,
		crc:        b.crc,
		depth:      b.depth + 1,
		parent:     b,
		tappedCard: 0,
		bookmark:   b.bookmark,
		moves:      b.moves,
		fmoves:     b.fmoves,
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

	child.setupPilesToCheck()
	child.findTapTargets()
	// child.tappedCard will be added later
	return child
}
*/

/*
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
*/

/*
func solve(root *Baize, b0 *Baize, depth int) {
	if b0.depth > depth || b0.percent >= 100 {
		println("limit reached", b0.depth, b0.percent)
		return
	}
	if b0.moves == 0 {
		println("stuck")
		return
	}
	println("b0 depth:", b0.depth, "moves:", b0.moves)
	println("b0 has", countTargets(b0), "tap targets")

	root.crc = root.calcCRC()

	// for each card with tap targets
	for ip0, p0 := range b0.piles {
		for ic0, c0 := range p0.cards {
			for nTapTarget := range c0.tapTargets {
				b1 := b0.createChild()
				p1 := b1.piles[ip0]
				c1 := p1.cards[ic0]
				if c0.id != c1.id {
					println("oops - cards are not the same")
				}
				tail := p1.makeTail(c1)
				oldCRC := b1.calcCRC()
				// move tail to a pre-determined destination
				b1.script.TailTapped(tail, nTapTarget)
				if b1.calcCRC() == oldCRC {
					println(ip0, ic0, "nothing changed!")
				} else {
					c0.tapTargets[nTapTarget].baize = b1
					// println("baize changed")
					// the tap caused the baize to change
					b1.script.AfterMove()
					b1.percent = b1.percentComplete()
					b1.tappedCard = c1.id
					// solve(root, b1, depth)
				}
			}
		}
	}

}
*/

type tapNode struct {
	cid      cardid.CardID
	percent  int
	depth    int
	crc      uint32
	parent   *tapNode
	children []*tapNode
}

func (b *Baize) solve2(root *tapNode, tn *tapNode, maxDepth int) {
	if tn.depth > maxDepth {
		// println("max depth reached")
		return
	}
	if b.moves == 0 {
		// println("no moves")
		return
	}
	// create a list of all the tap cards in this baize
	// because cards will be destroyed/recreated by undo
	var tapTargets []cardid.CardID = []cardid.CardID{}
	b.foreachCard(func(c *Card) {
		if c.tapTarget.dst != nil {
			tapTargets = append(tapTargets, c.id)
		}
	})
	// println(len(tapTargets), "tap targets found in baize")

	sb := b.newSavableBaize()
	oldPercent := b.percent

	for i := range tapTargets {
		c := b.findCard(tapTargets[i])
		tail := c.owner().makeTail(c)

		// make the move
		b.script.TailTapped(tail)
		b.script.AfterMove()
		b.findTapTargets()
		b.percent, _ = b.percentComplete()

		// we now have a new baize
		// create a child node with result of this move
		var node tapNode = tapNode{
			cid:      tapTargets[i],
			percent:  b.percent,
			depth:    tn.depth + 1,
			crc:      b.calcCRC(),
			parent:   tn,
			children: []*tapNode{},
		}
		if !findCRC(root, node.crc) {
			tn.children = append(tn.children, &node)
		} else {
			// println("skipping duplicate")
		}

		// go find children of this node
		b.solve2(root, &node, maxDepth)

		// revert baize to it's starting state
		b.updateFromSavable(sb)
		b.findTapTargets()
		b.percent = oldPercent
	}
}

func display(tn *tapNode) {
	log.Printf("depth %d complete %d%%", tn.depth, tn.percent)
	for _, tc := range tn.children {
		display(tc)
	}
}

func findCRC(tn *tapNode, crc uint32) bool {
	if tn.crc == crc {
		return true
	}
	for _, tc := range tn.children {
		if findCRC(tc, crc) {
			return true
		}
	}
	return false
}

func countNodes(tn *tapNode, pcount *int) {
	*pcount++
	for _, tc := range tn.children {
		countNodes(tc, pcount)
	}
}

func maxPercent(tn *tapNode, pmax *int, ptn **tapNode) {
	if tn.percent > *pmax {
		*pmax = tn.percent
		*ptn = tn
	}
	for _, tc := range tn.children {
		maxPercent(tc, pmax, ptn)
	}
}

// func (b *Baize) sanity() {
// 	for _, p := range b.piles {
// 		for _, c := range p.cards {
// 			if c.pile != p {
// 				println("insanity at", p.category, c.String())
// 			}
// 		}
// 	}
// }

func (b *Baize) Solve(maxDepth int) {

	var tn *tapNode = &tapNode{} // root node will be empty/dummy, except for children
	b.solve2(tn, tn, maxDepth)
	var count int
	countNodes(tn, &count)
	var max int
	var ptn **tapNode = &tn
	maxPercent(tn, &max, ptn)
	println("max depth", maxDepth, "nodes", count, "max percent", max, "card", (*ptn).cid.String())

	b.foreachCard(func(c *Card) {
		c.tapTarget.dst = nil
		c.tapTarget.weight = 0
	})
	// weight the next card to be tapped
	tn2 := *ptn
	var id cardid.CardID
	for tn2.parent != nil {
		id = tn2.cid
		if c, ok := b.cardMap[id.PackSuitOrdinal()]; ok {
			c.weight = 2
		}
		// println(tn2.cid.String())
		tn2 = tn2.parent
	}
	// display(tapTree)
	// b.findTapTargets()
}
