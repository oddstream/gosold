-- Mrs Mop

function BuildPiles()
	NewStock(-5, -5)
	for x = 0, 3 do
		NewDiscard(x, 0)
		NewDiscard(x+9, 0)
	end
	for x = 0, 12 do
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ANY)
		SetCompareFunction(t, "Append", "Down")
		SetCompareFunction(t, "Move", "DownSuit")
	end
end

function StartGame()
	-- 13 piles of 8 cards each
	local stock = Stock()
	for _, t in ipairs(Tableaux()) do
		for _ = 1, 8 do
			MoveCard(stock, t)
		end
	end
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Mrs._Mop"
end

function CardColors()
	return 4
end

function Packs()
	return 2
end