package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

type Spiderette struct {
	scriptBase
}

func (self *Spiderette) BuildPiles() {

	if self.cardColors == 0 {
		self.cardColors = 4
	}

	self.stock = self.baize.NewStock(newPileSlot(0, 0))

	for x := 3; x < 7; x++ {
		d := self.baize.NewDiscard(newPileSlot(x, 0))
		self.discards = append(self.discards, d)
	}

	for x := 0; x < 7; x++ {
		t := self.baize.NewTableau(newPileSlot(x, 1), FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.appendCmp2 = dyad.compare_Down
		t.moveCmp2 = dyad.compare_DownSuit
	}
}

func (self *Spiderette) StartGame() {
	var dealDown int = 1
	for _, pile := range self.tableaux {
		for i := 0; i < dealDown; i++ {
			if c := moveCard(self.stock, pile); c != nil {
				c.flipDown()
			}
		}
		dealDown++
		moveCard(self.stock, pile)
	}
	for _, pile := range self.tableaux {
		if c := pile.peek(); c != nil {
			c.flipUp()
		}
	}
	self.baize.setRecycles(0)
}

// default TailMoveError

// default TailAppendError

func (self *Spiderette) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].owner()
	switch pile.vtable.(type) {
	case *Stock:
		var tabCards, emptyTabs int
		for _, tab := range self.tableaux {
			if tab.Len() == 0 {
				emptyTabs++
			} else {
				tabCards += tab.Len()
			}
		}
		if emptyTabs > 0 && tabCards >= len(self.tableaux) {
			// TheGame.UI.ToastError("All empty tableaux must be filled before dealing a new row")
		} else {
			for _, tab := range self.tableaux {
				moveCard(self.stock, tab)
			}
		}
	default:
		tail[0].owner().vtable.tailTapped(tail)

	}
}

// func (*Spiderette) PileTapped(*Pile) {}
