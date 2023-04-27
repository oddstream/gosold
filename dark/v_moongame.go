package dark

//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type MoonGame struct {
	scriptBase
}

// SetBaize() not passed to Lua

// Reset() not passed to Lua

func (self *MoonGame) BuildPiles() {
	glob := self.baize.L.GetGlobal("BuildPiles")
	if fn, ok := glob.(*lua.LFunction); ok {
		err := self.baize.L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}, &lua.LUserData{Value: self})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("BuildPiles is not a Lua function, and there is no default")
	}
}

func (self *MoonGame) StartGame() {
	glob := self.baize.L.GetGlobal("StartGame") // glob == lua.LNil if it doesn't exist
	if fn, ok := glob.(*lua.LFunction); ok {
		err := self.baize.L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}, &lua.LUserData{Value: self})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("StartGame is not a Lua function, and there is no default")
	}
}

func (self *MoonGame) AfterMove() {
	glob := self.baize.L.GetGlobal("AfterMove") // glob == lua.LNil if it doesn't exist
	if fn, ok := glob.(*lua.LFunction); ok {
		err := self.baize.L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}, &lua.LUserData{Value: self})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		self.scriptBase.AfterMove()
	}
}

func (self *MoonGame) TailMoveError(tail []*Card) (bool, error) {
	var returnOk bool = true
	var returnErr error

	glob := self.baize.L.GetGlobal("TailMoveError") // glob == lua.LNil if it doesn't exist
	if fn, ok := glob.(*lua.LFunction); ok {
		err := self.baize.L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    2,
			Protect: true,
		}, &lua.LUserData{Value: self}, &lua.LUserData{Value: tail})
		if err != nil {
			fmt.Println(err)
		} else {
			// returnOk = self.getBool(1)
			returnOk = self.baize.L.CheckBool(1)
			returnErr = self.getError(2)
			self.baize.L.Pop(2)
		}
	} else {
		returnOk, returnErr = self.scriptBase.TailMoveError(tail)
	}
	return returnOk, returnErr
}

func (self *MoonGame) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	var returnOk bool = true
	var returnErr error

	glob := self.baize.L.GetGlobal("TailAppendError") // glob == lua.LNil if it doesn't exist
	if fn, ok := glob.(*lua.LFunction); ok {
		// udMoonGame := self.baize.L.NewUserData()
		// udMoonGame.Value = self
		// udPile := self.baize.L.NewUserData()
		// udPile.Value = dst
		// udTail := self.baize.L.NewUserData()
		// udTail.Value = tail
		err := self.baize.L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    2,
			Protect: true,
		}, &lua.LUserData{Value: self}, &lua.LUserData{Value: dst}, &lua.LUserData{Value: tail})
		// }, udMoonGame, udPile, udTail)
		if err != nil {
			fmt.Println(err)
		} else {
			// returnOk = self.getBool(1)
			returnOk = self.baize.L.CheckBool(1)
			returnErr = self.getError(2)
			self.baize.L.Pop(2)
		}
	} else {
		returnOk, returnErr = self.scriptBase.TailAppendError(dst, tail)
	}
	return returnOk, returnErr
}

func (self *MoonGame) TailTapped(tail []*Card) {
	glob := self.baize.L.GetGlobal("TailTapped") // glob == lua.LNil if it doesn't exist
	if fn, ok := glob.(*lua.LFunction); ok {
		err := self.baize.L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}, &lua.LUserData{Value: self}, &lua.LUserData{Value: tail})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		self.scriptBase.TailTapped(tail)
	}
}

func (self *MoonGame) PileTapped(pile *Pile) {
	glob := self.baize.L.GetGlobal("PileTapped") // glob == lua.LNil if it doesn't exist
	if fn, ok := glob.(*lua.LFunction); ok {
		err := self.baize.L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}, &lua.LUserData{Value: self}, &lua.LUserData{Value: pile})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		self.scriptBase.PileTapped(pile)
	}
}

// TODO Cells(), Discards(), Foundations(), Reserves(), Stock(), Tableaux(), Wastes()
// return a Lua table of *Pile

func (self *MoonGame) Complete() bool {
	var complete bool

	glob := self.baize.L.GetGlobal("Complete") // glob == lua.LNil if it doesn't exist
	if fn, ok := glob.(*lua.LFunction); ok {
		err := self.baize.L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    1,
			Protect: true,
		}, &lua.LUserData{Value: self})
		if err != nil {
			fmt.Println(err)
		} else {
			// complete = self.getBool(1)
			complete = self.baize.L.CheckBool(1)
			self.baize.L.Pop(1)
		}
	} else {
		complete = self.scriptBase.Complete()
	}
	return complete
}

func (self *MoonGame) Wikipedia() string {
	var str string

	glob := self.baize.L.GetGlobal("Wikipedia") // glob == lua.LNil if it doesn't exist
	if fn, ok := glob.(*lua.LFunction); ok {
		err := self.baize.L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    1,
			Protect: true,
		}, &lua.LUserData{Value: self})
		if err != nil {
			fmt.Println(err)
		} else {
			str = self.baize.L.CheckString(1)
			self.baize.L.Pop(1)
		}
	} else {
		str = self.scriptBase.Wikipedia()
	}
	return str
}

// functions called by Lua to do DARK things

// moonPileList returns a Lua list of piles of a given category
func _moonGetPiles(L *lua.LState, category string) int {
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		tab := L.NewTable()
		var piles []*Pile
		switch category {
		case "Cells":
			piles = moonGame.cells
		case "Discards":
			piles = moonGame.discards
		case "Foundations":
			piles = moonGame.foundations
		case "Reserves":
			piles = moonGame.reserves
		case "Tableaux":
			piles = moonGame.tableaux
		case "Wastes":
			piles = moonGame.wastes
		default:
			fmt.Println("Cannot get", category, "piles")
		}
		for i, p := range piles {
			udc := L.NewUserData()
			udc.Value = p
			L.RawSetInt(tab, i+1, udc)
		}
		L.Push(tab)
		return 1
	}
	return 0
}

func moonNewCell(L *lua.LState) int {
	var pile *Pile
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		x := L.ToNumber(2)
		y := L.ToNumber(3)
		baize := moonGame.baize
		pile = baize.NewCell(PileSlot{X: float32(x), Y: float32(y), Deg: 0})
		moonGame.cells = append(moonGame.cells, pile)
	}
	udPile := L.NewUserData()
	udPile.Value = pile
	L.Push(udPile)
	return 1
}

func moonGetCells(L *lua.LState) int {
	return _moonGetPiles(L, "Cells")
}

func moonNewFoundation(L *lua.LState) int {
	var pile *Pile
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		x := L.ToNumber(2)
		y := L.ToNumber(3)
		baize := moonGame.baize
		pile = baize.NewFoundation(PileSlot{X: float32(x), Y: float32(y), Deg: 0})
		moonGame.foundations = append(moonGame.foundations, pile)
	}
	udPile := L.NewUserData()
	udPile.Value = pile
	L.Push(udPile)
	// L.Push(&lua.LUserData{Value: pile})
	return 1
}

func moonGetFoundations(L *lua.LState) int {
	return _moonGetPiles(L, "Foundations")
}

func moonNewReserve(L *lua.LState) int {
	var pile *Pile
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		x := L.ToNumber(2)
		y := L.ToNumber(3)
		fanType := FanType(L.CheckInt(4))
		baize := moonGame.baize
		pile = baize.NewReserve(PileSlot{X: float32(x), Y: float32(y), Deg: 0}, fanType)
		moonGame.reserves = append(moonGame.reserves, pile)
	}
	udPile := L.NewUserData()
	udPile.Value = pile
	L.Push(udPile)
	return 1
}

func moonGetReserves(L *lua.LState) int {
	return _moonGetPiles(L, "Reserves")
}

func moonNewTableau(L *lua.LState) int {
	var pile *Pile
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		x := L.ToNumber(2)
		y := L.ToNumber(3)
		fanType := FanType(L.CheckInt(4))
		moveType := MoveType(L.CheckInt(5))
		baize := moonGame.baize
		pile = baize.NewTableau(PileSlot{X: float32(x), Y: float32(y), Deg: 0}, fanType, moveType)
		moonGame.tableaux = append(moonGame.tableaux, pile)
	}
	udPile := L.NewUserData()
	udPile.Value = pile
	L.Push(udPile)
	return 1
}

func moonGetTableaux(L *lua.LState) int {
	return _moonGetPiles(L, "Tableaux")
}

func moonNewWaste(L *lua.LState) int {
	var pile *Pile
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		x := L.ToNumber(2)
		y := L.ToNumber(3)
		fanType := FanType(L.CheckInt(4))
		baize := moonGame.baize
		pile = baize.NewWaste(PileSlot{X: float32(x), Y: float32(y), Deg: 0}, fanType)
		moonGame.wastes = append(moonGame.wastes, pile)
	}
	udPile := L.NewUserData()
	udPile.Value = pile
	L.Push(udPile)
	return 1
}

func moonGetWaste(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		if len(moonGame.wastes) > 0 {
			udPile := L.NewUserData()
			udPile.Value = moonGame.wastes[0]
			L.Push(udPile)
		}
	}
	return 1
}

func moonGetWastes(L *lua.LState) int {
	return _moonGetPiles(L, "Wastes")
}

func moonNewStock(L *lua.LState) int {
	var pile *Pile
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		x := L.ToNumber(2)
		y := L.ToNumber(3)
		baize := moonGame.baize
		pile = baize.NewStock(PileSlot{X: float32(x), Y: float32(y), Deg: 0})
		moonGame.stock = pile
	}
	udPile := L.NewUserData()
	udPile.Value = pile
	L.Push(udPile)
	return 1
}

func moonGetStock(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		udPile := L.NewUserData()
		udPile.Value = moonGame.stock
		L.Push(udPile)
	}
	return 1
}

func moonMoveCard(L *lua.LState) int {
	// TODO check arity with GetTop() == 3
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udSrc := L.CheckUserData(2)
		if src, ok := udSrc.Value.(*Pile); ok {
			udDst := L.CheckUserData(3)
			if dst, ok := udDst.Value.(*Pile); ok {
				card := moveCard(src, dst)
				udCard := L.NewUserData()
				udCard.Value = card
				L.Push(udCard)
				return 1
			}
		}
	}
	return 0
}

func moonFlipDown(L *lua.LState) int {
	// TODO check arity with GetTop() == 3
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udCard := L.CheckUserData(2)
		if card, ok := udCard.Value.(*Card); ok {
			card.flipDown()
		}
	}
	return 0
}

func moonFlipUp(L *lua.LState) int {
	// TODO check arity with GetTop() == 2
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udCard := L.CheckUserData(2)
		if card, ok := udCard.Value.(*Card); ok {
			card.flipUp()
		}
	}
	return 0
}

func moonOrdinal(L *lua.LState) int {
	// TODO check arity with GetTop() == 2
	var ord int
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udCard := L.CheckUserData(2)
		if card, ok := udCard.Value.(*Card); ok {
			ord = card.Ordinal()
		}
	}
	L.Push(lua.LNumber(ord))
	return 1
}

func moonOwner(L *lua.LState) int {
	// TODO check arity with GetTop() == 2
	var pile *Pile
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udCard := L.CheckUserData(2)
		if card, ok := udCard.Value.(*Card); ok {
			pile = card.owner()
		}
	}
	ud = L.NewUserData()
	ud.Value = pile
	L.Push(ud)
	return 1
}

func moonFirst(L *lua.LState) int {
	// TODO check arity with GetTop() == 2
	var card *Card
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.CheckUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			card = pile.cards[0]
		} else if tail, ok := udPile.Value.([]*Card); ok {
			card = tail[0]
		}
	}
	udCard := L.NewUserData()
	udCard.Value = card
	L.Push(udCard)
	return 1
}
func moonPeek(L *lua.LState) int {
	// TODO check arity with GetTop() == 2
	var card *Card
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.CheckUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			card = pile.peek()
		} else if tail, ok := udPile.Value.([]*Card); ok {
			card = tail[len(tail)-1]
		}
	}
	udCard := L.NewUserData()
	udCard.Value = card
	L.Push(udCard)
	return 1
}

func moonLen(L *lua.LState) int {
	// TODO check arity with GetTop() == 2
	var length int
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.CheckUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			length = pile.Len()
		} else if tail, ok := udPile.Value.([]*Card); ok {
			length = len(tail)
		}
	}
	L.Push(lua.LNumber(length)) // LNumber is a type alias for float64
	return 1
}

func moonCategory(L *lua.LState) int {
	// TODO check arity with GetTop() == 2
	var cat string
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.CheckUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			cat = pile.category
		}
	}
	L.Push(lua.LString(cat))
	return 1
}

var moonCompareFunctions = map[string]dyadCmpFunc{
	"Any":              dyad.compare_Any,
	"Up":               dyad.compare_Up,
	"UpWrap":           dyad.compare_UpWrap,
	"Down":             dyad.compare_Down,
	"DownWrap":         dyad.compare_DownWrap,
	"UpOrDown":         dyad.compare_UpOrDown,
	"UpOrDownWrap":     dyad.compare_UpOrDownWrap,
	"Color":            dyad.compare_Color,
	"AltColor":         dyad.compare_AltColor,
	"Suit":             dyad.compare_Suit,
	"OtherSuit":        dyad.compare_OtherSuit,
	"DownColor":        dyad.compare_DownColor,
	"DownAltColor":     dyad.compare_DownAltColor,
	"DownAltColorWrap": dyad.compare_DownAltColorWrap,
	"UpAltColor":       dyad.compare_UpAltColor,
	"UpSuit":           dyad.compare_UpSuit,
	"DownSuit":         dyad.compare_DownSuit,
	"UpOrDownSuit":     dyad.compare_UpOrDownSuit,
	"UpOrDownSuitWrap": dyad.compare_UpOrDownSuitWrap,
	"DownOtherSuit":    dyad.compare_DownOtherSuit,
	"UpSuitWrap":       dyad.compare_UpSuitWrap,
	"DownSuitWrap":     dyad.compare_DownSuitWrap,
}

func moonSetCompareFunction(L *lua.LState) int {
	// TODO check arity with GetTop() == 4
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.CheckUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			typ := L.CheckString(3)
			fn := L.CheckString(4)
			// TODO check typ and fn
			switch typ {
			case "Append":
				if pile.appendCmp2, ok = moonCompareFunctions[string(fn)]; !ok {
					fmt.Println("Unknown append compare function", string(fn))
				}
			case "Move":
				if pile.moveCmp2, ok = moonCompareFunctions[string(fn)]; !ok {
					fmt.Println("Unknown move compare function", string(fn))
				}
			default:
				fmt.Println("Unknown type", typ)
			}
		}
	}
	return 0
}

func moonLabel(L *lua.LState) int {
	// TODO check arity with GetTop() == 2
	var str string
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.CheckUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			str = pile.Label()
		} else {
			fmt.Println("SetLabel expecting a *Pile, got a", ud.Type().String())
		}
	} else {
		fmt.Println("SetLabel expecting a MoonGame, got a", ud.Type().String())
	}
	L.Push(lua.LString(str))
	return 1
}

func moonSetLabel(L *lua.LState) int {
	// TODO check arity with GetTop() == 3
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.CheckUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			str := L.CheckString(3)
			pile.setLabel(str)
		} else {
			fmt.Println("SetLabel expecting a *Pile, got a", ud.Type().String())
		}
	} else {
		fmt.Println("SetLabel expecting a MoonGame, got a", ud.Type().String())
	}
	return 0
}

func moonSetRecycles(L *lua.LState) int {
	// TODO check arity with GetTop() == 2
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		n := L.CheckInt(2)
		moonGame.baize.setRecycles(n)
	} else {
		fmt.Println("SetLabel expecting a MoonGame, got a", ud.Type().String())
	}
	return 0
}

func moonCompareEmpty(L *lua.LState) int {
	var result bool
	var err error
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.CheckUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			udTail := L.CheckUserData(3)
			if tail, ok := udTail.Value.([]*Card); ok {
				result, err = compare_Empty(pile, tail)
			} else {
				fmt.Println("CompareEmpty arg 3 is not a tail, got a", udTail.Type().String())
			}
		} else {
			fmt.Println("CompareEmpty arg 2 is not a pile, got a", udPile.Type().String())
		}
	}
	L.Push(lua.LBool(result))
	if err == nil {
		L.Push(lua.LString(""))
	} else {
		L.Push(lua.LString(err.Error()))
	}
	return 2 // bool, error string
}

func moonCompareAppend(L *lua.LState) int {
	var result bool
	var err error
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.CheckUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			udTail := L.CheckUserData(3)
			if tail, ok := udTail.Value.([]*Card); ok {
				result, err = pile.appendCmp2(dyad{pile.peek(), tail[0]})
			} else {
				fmt.Println("CompareAppend arg 3 is not a tail, got a", udTail.Type().String())
			}
		} else {
			fmt.Println("CompareAppend arg 2 is not a pile, got a", udPile.Type().String())
		}
	}
	L.Push(lua.LBool(result))
	if err == nil {
		L.Push(lua.LString(""))
	} else {
		L.Push(lua.LString(err.Error()))
	}
	return 2 // bool, error string
}

func moonCompareMove(L *lua.LState) int {
	var result bool
	var err error
	ud := L.CheckUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.CheckUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			udTail := L.CheckUserData(3)
			if tail, ok := udTail.Value.([]*Card); ok {
				result, err = pile.moveCmp2(dyad{pile.peek(), tail[0]})
			} else {
				fmt.Println("CompareMove arg 3 is not a tail, got a", udTail.Type().String())
			}
		} else {
			fmt.Println("CompareMove arg 2 is not a pile, got a", udPile.Type().String())
		}
	}
	L.Push(lua.LBool(result))
	if err == nil {
		L.Push(lua.LString(""))
	} else {
		L.Push(lua.LString(err.Error()))
	}
	return 2 // bool, error string
}

func moonToast(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		str := L.CheckString(2)
		moonGame.baize.fnNotify(MessageEvent, str)
	}
	return 0
}

// register functions with Lua

func registerMoonFunctions(L *lua.LState) {
	// Any function registered with GopherLua is a lua.LGFunction, defined in value.go
	// type LGFunction func(*LState) int
	funcs := []struct {
		name string
		fn   lua.LGFunction //func(*lua.LState) int
	}{
		// Pile creation
		{"NewCell", moonNewCell},
		{"GetCells", moonGetCells},
		// {"NewDiscard", moonNewDiscard},
		// {"GetDiscards", moonGetDiscards},
		{"NewFoundation", moonNewFoundation},
		{"GetFoundations", moonGetFoundations},
		{"NewReserve", moonNewReserve},
		{"GetReserves", moonGetReserves},
		{"NewStock", moonNewStock},
		{"GetStock", moonGetStock},
		{"NewTableau", moonNewTableau},
		{"GetTableaux", moonGetTableaux},
		{"NewWaste", moonNewWaste},
		{"GetWaste", moonGetWaste},
		{"GetWastes", moonGetWastes},

		// Baize
		{"MoveCard", moonMoveCard},
		{"SetRecycles", moonSetRecycles},

		// Pile
		{"Category", moonCategory},
		{"Label", moonLabel},
		{"Len", moonLen},     // pile or tail
		{"First", moonFirst}, // pile or tail
		{"Peek", moonPeek},   // pile or tail
		{"SetCompareFunction", moonSetCompareFunction},
		{"SetLabel", moonSetLabel},

		// Card
		{"FlipDown", moonFlipDown},
		{"FlipUp", moonFlipUp},
		{"Ordinal", moonOrdinal},
		{"Owner", moonOwner},

		// Other
		{"CompareEmpty", moonCompareEmpty},
		{"CompareAppend", moonCompareAppend},
		{"CompareMove", moonCompareMove},
		{"Toast", moonToast},
	}
	for _, f := range funcs {
		L.SetGlobal(f.name, L.NewFunction(f.fn))
	}

	L.SetGlobal("FAN_NONE", lua.LNumber(FAN_NONE))
	L.SetGlobal("FAN_DOWN", lua.LNumber(FAN_DOWN))
	L.SetGlobal("FAN_LEFT", lua.LNumber(FAN_LEFT))
	L.SetGlobal("FAN_RIGHT", lua.LNumber(FAN_RIGHT))
	L.SetGlobal("FAN_DOWN3", lua.LNumber(FAN_DOWN3))
	L.SetGlobal("FAN_LEFT3", lua.LNumber(FAN_LEFT3))
	L.SetGlobal("FAN_RIGHT3", lua.LNumber(FAN_RIGHT3))

	L.SetGlobal("MOVE_NONE", lua.LNumber(MOVE_NONE))
	L.SetGlobal("MOVE_ANY", lua.LNumber(MOVE_ANY))
	L.SetGlobal("MOVE_ONE", lua.LNumber(MOVE_ONE))
	L.SetGlobal("MOVE_ONE_PLUS", lua.LNumber(MOVE_ONE_PLUS))
	L.SetGlobal("MOVE_ONE_OR_ALL", lua.LNumber(MOVE_ONE_OR_ALL))
}

// utility functions

func (self *MoonGame) getBool(stackPos int) bool {
	val := self.baize.L.Get(stackPos)
	if _, ok := val.(lua.LBool); !ok {
		fmt.Println("Lua function returned a", val.Type().String(), "instead of a bool")
	} else {
		if val == lua.LTrue {
			return true
		} else {
			return false
		}
	}
	return false
}

func (self *MoonGame) getString(stackPos int) string {
	val := self.baize.L.Get(stackPos)
	if str, ok := val.(lua.LString); !ok {
		fmt.Println("Lua function returned a", val.Type().String(), "instead of a string")
	} else {
		return string(str)
	}
	return ""
}

func (self *MoonGame) getNumber(stackPos int) float64 {
	val := self.baize.L.Get(stackPos)
	if n, ok := val.(lua.LNumber); !ok {
		fmt.Println("Lua function returned a", val.Type().String(), "instead of a number")
	} else {
		return float64(n)
	}
	return 0
}

func (self *MoonGame) getError(stackPos int) error {
	val := self.baize.L.Get(stackPos)
	if str, ok := val.(lua.LString); !ok {
		fmt.Println("Lua function returned a", val.Type().String(), "instead of a string")
	} else {
		if string(str) != "" {
			return errors.New(string(str))
		}
	}
	return nil
}
