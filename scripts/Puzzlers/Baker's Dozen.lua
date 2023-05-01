-- Baker's Dozen

function BuildPiles()
	NewStock(-5,-5)
	for x = 0, 6 do
		local t = NewTableau(x, 0, FAN_DOWN, MOVE_ONE)
		SetCompareFunction(t, "Append", "Down")
		SetLabel(t, "x")
	end
	for x = 0, 5 do
		local t = NewTableau(x + 0.5, 3, FAN_DOWN, MOVE_ONE)
		SetCompareFunction(t, "Append", "Down")
		SetLabel(t, "x")
	end
	for y = 0, 3 do
		local f = NewFoundation(8, y)
		SetCompareFunction(f, "Append", "UpSuit")
		SetLabel(f, "A")
	end
--[[
	for x := 0; x < 6; x++ {
		// stock is pile index 0
		// tableaux are piles index 1 .. 13
		self.tableaux[x].boundary = 1 + x + 7
	}
	self.tableaux[6].boundary = 1 + 12
]]
end

function StartGame()
	for _, t in ipairs(Tableaux()) do
		for _ = 1,4 do
			MoveCard(Stock(), t)
		end
		Bury(t, 13)
	end
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Baker%27s_Dozen_(card_game)"
end

function CardColors()
	return 1
end
