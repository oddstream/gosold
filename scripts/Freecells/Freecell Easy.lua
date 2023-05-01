-- Freecell Easy

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
	local stock = Stock()
	local founds = Foundations()

	local c
	c = Extract(stock, 0, 1, CLUB)
	Push(founds[1], c)
	c = Extract(stock, 0, 1, DIAMOND)
	Push(founds[2], c)
	c = Extract(stock, 0, 1, HEART)
	Push(founds[3], c)
	c = Extract(stock, 0, 1, SPADE)
	Push(founds[4], c)

	for _, t in ipairs(Tableaux()) do
		for _ = 1, 6 do
			MoveCard(stock, t)
		end
		Bury(t, 13)
	end
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Freecell"
end