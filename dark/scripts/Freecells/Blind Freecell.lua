-- Blind Freecell

function BuildPiles()
	NewStock(-5, -5)
	for x = 0, 3 do
		NewCell(x, 0)
	end
	for x = 4, 7 do
		local f = NewFoundation(x, 0)
		SetCompareFunction(f, "Append", "UpSuit")
		SetLabel(f, "A")
	end
	for x = 0, 7 do
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ONE_PLUS)
		SetCompareFunction(t, "Append", "DownAltColor")
		SetCompareFunction(t, "Move", "DownAltColor")
	end

end

function StartGame()
	-- 4 tabs [0 .. 3] with 7 cards
	-- 4 tabs [4 .. 7] with 6 cards
	local stock = Stock()
	local tabs = Tableaux()
	for i = 1, 4 do
		for _ = 1, 7 do
			local c = MoveCard(stock, tabs[i])
			FlipDown(c)
		end
	end
	for i = 5, 8 do
		for _ = 1, 6 do
			local c = MoveCard(stock, tabs[i])
			FlipDown(c)
		end
	end
	for _, t in ipairs(tabs) do
		local c = Peek(t)
		FlipUp(c)
	end
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/FreeCell"
end
