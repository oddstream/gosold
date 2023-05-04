-- Martha

function BuildPiles()
	NewStock(-5, -5)
	for x = 8, 11 do
		local f = NewFoundation(x, 0)
		SetCompareFunction(f, "Append", "UpSuit")
		SetLabel(f, "A")
	end
	for x = 0, 11 do
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ANY)
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
		c = MoveCard(stock, t)
		FlipDown(c)
		MoveCard(stock, t)	-- face up
		c = MoveCard(stock, t)
		FlipDown(c)
		MoveCard(stock, t)	-- face up
	end
end

-- One card can be moved at a time, but sequences that are already built can be moved, in part or in whole, as unit.
-- But when a gap occurs, it can be filled only with a single card.
function TailAppendError(dst, tail)
	if Len(dst) == 0 and Category(dst) == "Tableau" then
		if Len(tail) ~= 1 then
			return false, "Empty Tableau piles can only be filled by a single card"
		end
	end
	return DefaultTailAppendError(dst, tail)
	-- if Len(dst) == 0 then
	-- 	return CompareEmpty(dst, tail)
	-- end
	-- return CompareAppend(dst, tail)
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Martha_(solitaire)"
end
