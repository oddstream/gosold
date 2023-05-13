-- Baker's Dozen

function BuildPiles()
	NewStock(-5,-5)

	for x = 0, 6 do
		local t = NewTableau(x, 0, FAN_DOWN, MOVE_ONE)
		SetCompareFunction(t, "Append", "Down")
		SetLabel(t, "X")
	end
	for x = 0, 5 do
		local t = NewTableau(x + 0.5, 3, FAN_DOWN, MOVE_ONE)
		SetCompareFunction(t, "Append", "Down")
		SetLabel(t, "X")
	end
	local ts = Tableaux()
	SetBoundary(ts[1], 8)
	SetBoundary(ts[2], 9)
	SetBoundary(ts[3], 10)
	SetBoundary(ts[4], 11)
	SetBoundary(ts[5], 12)
	SetBoundary(ts[6], 13)
	SetBoundary(ts[7], 13)

	for y = 0, 3 do
		local f = NewFoundation(8, y)
		SetCompareFunction(f, "Append", "UpSuit")
		SetLabel(f, "A")
	end
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
