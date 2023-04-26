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
	glob := self.baize.L.GetGlobal("BuildPiles") // glob == lua.LNil if it doesn't exist
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
		err := self.baize.L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    2,
			Protect: true,
		}, &lua.LUserData{Value: self}, &lua.LUserData{Value: dst}, &lua.LUserData{Value: tail})
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
	glob := self.baize.L.GetGlobal("TailTapped") // glob == lua.LNil if it doesn't exist
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

// Cells(), Discards(), Foundations(), Reserves(), Stock(), Tableaux(), Wastes() not used in Lua (?)

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

func moonNewCell(L *lua.LState) int {
	var pile *Pile
	ud := L.ToUserData(1)
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

func moonNewFoundation(L *lua.LState) int {
	var pile *Pile
	ud := L.ToUserData(1)
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
	return 1
}

func moonNewTableau(L *lua.LState) int {
	var pile *Pile
	ud := L.ToUserData(1)
	if moonGame, ok := ud.Value.(*MoonGame); ok {
		x := L.ToNumber(2)
		y := L.ToNumber(3)
		moveType := MoveType(L.ToNumber(4))
		fanType := FanType(L.ToNumber(5))
		baize := moonGame.baize
		pile = baize.NewTableau(PileSlot{X: float32(x), Y: float32(y), Deg: 0}, fanType, moveType)
		moonGame.tableaux = append(moonGame.tableaux, pile)
	}
	udPile := L.NewUserData()
	udPile.Value = pile
	L.Push(udPile)
	return 1
}

func moonNewStock(L *lua.LState) int {
	var pile *Pile
	ud := L.ToUserData(1)
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

func moonMoveCard(L *lua.LState) int {
	// TODO check arity with GetTop() == 3
	ud := L.ToUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udSrc := L.ToUserData(2)
		if src, ok := udSrc.Value.(*Pile); ok {
			udDst := L.ToUserData(3)
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
				pile.appendCmp2 = moonCompareFunctions[string(fn)]
			case "Move":
				pile.moveCmp2 = moonCompareFunctions[string(fn)]
			}
		}
	}
	return 0
}

func moonSetLabel(L *lua.LState) int {
	// TODO check arity with GetTop() == 3
	ud := L.ToUserData(1)
	if _, ok := ud.Value.(*MoonGame); ok {
		udPile := L.ToUserData(2)
		if pile, ok := udPile.Value.(*Pile); ok {
			str := L.ToString(3)
			pile.setLabel(str)
		} else {
			fmt.Println("SetLabel expecting a *Pile, got a", ud.Type().String())
		}
	} else {
		fmt.Println("SetLabel expecting a MoonGame, got a", ud.Type().String())
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
		{"NewCell", moonNewCell},
		// {"NewDiscard", moonNewDiscard},
		{"NewFoundation", moonNewFoundation},
		// {"NewReserve", moonNewReserve},
		{"NewStock", moonNewStock},
		{"NewTableau", moonNewTableau},
		// {"NewWaste", moonNewWaste},
		{"MoveCard", moonMoveCard},
		{"NewStock", moonNewStock},
		{"SetCompareFunction", moonSetCompareFunction},
		{"SetLabel", moonSetLabel},
	}
	for _, f := range funcs {
		L.SetGlobal(f.name, L.NewFunction(f.fn))
	}

	L.SetGlobal("FAN_NONE", lua.LNumber(FAN_NONE))
	L.SetGlobal("FAN_DOWN", lua.LNumber(FAN_DOWN))
	L.SetGlobal("FAN_LEFT", lua.LNumber(FAN_LEFT))
	L.SetGlobal("FAN_RIGHT", lua.LNumber(FAN_RIGHT))
	L.SetGlobal("FAN_DOWN3", lua.LNumber(FAN_DOWN))
	L.SetGlobal("FAN_LEFT3", lua.LNumber(FAN_LEFT))
	L.SetGlobal("FAN_RIGHT3", lua.LNumber(FAN_RIGHT))

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
