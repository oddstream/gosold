-- Simple Simon

function BuildPiles()
	NewStock(-5, -5)
	for x = 3, 6 do
		NewDiscard(x, 0)
	end
	for x = 0, 9 do
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ANY)
		SetCompareFunction(t, "Append", "Down")
		SetCompareFunction(t, "Move", "DownSuit")
	end
end

function StartGame()
	-- 3 piles of 8 cards each
	local stock = Stock()
	local tx = Tableaux()
	for i = 1, 3 do
		local t = tx[i]
		for _ = 1, 8 do
			MoveCard(stock, t)
		end
	end
	local deal = 7
	for i = 4, 10 do
		local t = tx[i]
		for _ = 1, deal do
			MoveCard(stock, t)
		end
		deal = deal - 1
	end
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Simple_Simon_(solitaire)"
end

function CardColors()
	return 4
end