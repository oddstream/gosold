-- Easthaven

function BuildPiles()
	NewStock(0, 0)
	for x = 3, 6 do
		local f = NewFoundation(x, 0)
		SetCompareFunction(f, "Append", "UpSuit")
		SetLabel(f, "A")
	end
	for x = 0, 6 do
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ANY)
		SetCompareFunction(t, "Append", "DownAltColor")
		SetCompareFunction(t, "Move", "DownAltColor")
		SetLabel(t, "K")
	end
end

function StartGame()
	local stock = Stock()
	for _, t in ipairs(Tableaux()) do
		for _ = 1, 2 do
			local c = MoveCard(stock, t)
			FlipDown(c)
		end
		MoveCard(stock, t)
	end
	SetRecycles(0)
end

function TailTapped(tail)
	if Category(Owner(First(tail))) == "Stock" then
		for _, t in ipairs(Tableaux()) do
			MoveCard(Stock(), t)
		end
	else
		DefaultTailTapped(tail)
	end
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Klondike_(solitaire)"
end
