-- Eight Off

function BuildPiles()
	NewStock(-5, -5)
	for x = 0, 7 do
		NewCell(x, 0)
	end
	for y = 0, 3 do
		local f = NewFoundation(9, y)
		SetCompareFunction(f, "Append", "UpSuit")
		SetLabel(f, "A")
	end
	for x = 0, 7 do
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ONE_PLUS)
		SetCompareFunction(t, "Append", "DownSuit")
		SetCompareFunction(t, "Move", "DownSuit")
		SetLabel(t, "K")
	end

end

function StartGame()
	local stock = Stock()
	local cells = Cells()
	for i = 1, 4 do
		MoveCard(stock, cells[i])
	end
	for _, t in ipairs(Tableaux()) do
		for _ = 1, 6 do
			MoveCard(stock, t)
		end
	end
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Eight_Off"
end

function CardColors()
	return 4
end
