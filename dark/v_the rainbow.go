package dark

import (
	"errors"

	"oddstream.games/gosold/util"
)

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

type TheRainbow struct {
	scriptBase
	tab, rainbow []*Pile
}

var rainbow_slots []PileSlot = []PileSlot{
	{0, 3, -90},
	{0.1 * 1.3333, 3.0 - 0.78, -75},
	{0.4 * 1.3333, 3.0 - 1.5, -60},
	{.88 * 1.3333, 3.0 - 2.12, -45},
	{1.5 * 1.3333, 3.0 - 2.6, -30},
	{2.22 * 1.3333, 3.0 - 2.9, -15},
	{3.0 * 1.3333, 0.0, 0},
	{3.78 * 1.3333, 3.0 - 2.9, 15},
	{4.5 * 1.3333, 3.0 - 2.6, 30},
	{5.12 * 1.3333, 3.0 - 2.12, 45},
	{5.6 * 1.3333, 3.0 - 1.5, 60},
	{5.9 * 1.3333, 3.0 - 0.78, 75},
	{8, 3.0, 90},
}

// x, y = radius * cos angle, radius * sin angle
// center x,y is 7,3?

func (self *TheRainbow) BuildPiles() {
	/*
		for i := len(rainbow_slots) - 1; i >= 0; i-- {
			// angle := float64(rainbow_slots[i].Deg) * math.Pi / 180.0
			angle := float64(i*15) * math.Pi / 180.0
			rainbow_slots[i].X = float32(3.0*math.Cos(angle)) + 3.0
			rainbow_slots[i].Y = float32(3.0 * math.Sin(angle))
			log.Printf("%.2f, %.2f, %d", rainbow_slots[i].X, rainbow_slots[i].Y, i*15)
		}
	*/
	self.stock = self.baize.NewStock(newHiddenPileSlot())

	self.rainbow = nil
	for _, rs := range rainbow_slots {
		t := self.baize.NewTableau(rs, FAN_NONE, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		self.rainbow = append(self.rainbow, t)
		t.appendCmp2 = dyad.compare_Down
		// t.setLabel("X")
	}

	for x := float32(2.5); x < 6.0; x += 1.0 {
		f := self.baize.NewFoundation(PileSlot{x, 2, 0})
		self.foundations = append(self.foundations, f)
		f.appendCmp2 = dyad.compare_UpSuit
		f.setLabel("A")
	}

	self.tab = nil
	for x := 1; x < 8; x++ {
		t := self.baize.NewTableau(PileSlot{float32(x), 4.0, 0}, FAN_DOWN, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		self.tab = append(self.tab, t)
		t.appendCmp2 = dyad.compare_Down
	}

}

func (self *TheRainbow) StartGame() {

	for _, f := range self.foundations {
		if c := self.stock.extractOrdinal(1); c != nil {
			f.push(c)
		}
	}

	for _, t := range self.tab {
		for i := 1; i < 6; i++ {
			moveCard(self.stock, t)
		}
	}

	for _, t := range self.rainbow {
		moveCard(self.stock, t)
	}

	self.baize.setRecycles(0)
}

func (self *TheRainbow) AfterMove() {
}

/*
	"If you reach a point at which you can neither pack nor build any further,
	you have the right to move any one available card to the rainbow,
	if there is a higher card there upon which to pack it,
	as an eight, upon which you may pack a seven"
*/

func (*TheRainbow) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (self *TheRainbow) TailAppendError(dst *Pile, tail []*Card) (bool, error) {

	if dst.Empty() {
		if util.Contains(self.rainbow, dst) {
			return false, errors.New("Cannot put cards in an empty rainbow pile")
		}
		return compare_Empty(dst, tail[0])
	}

	src := tail[0].pile
	if util.Contains(self.rainbow, src) && util.Contains(self.rainbow, dst) {
		return false, errors.New("Cannot move cards between rainbow piles")
	}
	return dst.appendCmp2(dyad{dst.peek(), tail[0]})
}

func (self *TheRainbow) TailTapped(tail []*Card) {
	pile := tail[0].owner()
	pile.vtable.tailTapped(tail)
}

//func (self *TheRainbow) PileTapped(pile *Pile) {}
