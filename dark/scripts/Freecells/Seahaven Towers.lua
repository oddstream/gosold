-- Seahaven Towers

function BuildPiles()
	NewStock(-5, -5)
	for x = 0, 3 do
		NewCell(x, 0)
	end
	for x = 6, 9 do
		local f = NewFoundation(x, 0)
		SetCompareFunction(f, "Append", "UpSuit")
		SetLabel(f, "A")
	end
	for x = 0, 9 do
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ONE_PLUS)
		SetCompareFunction(t, "Append", "DownSuit")
		SetCompareFunction(t, "Move", "DownSuit")
		SetLabel(t, "K")
	end
end

function StartGame()
	local stock = Stock()
	local cells = Cells()
	for _, t in ipairs(Tableaux()) do
		for _ = 1, 5 do
			MoveCard(stock, t)
		end
	end
	MoveCard(stock, cells[1])
	MoveCard(stock, cells[2])
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Seahaven_Towers"
end

function CardColors()
	return 4
end
