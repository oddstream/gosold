-- Agnes Bernauer

local ordToChar = {"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

function BuildPiles()
	NewStock(0, 0)
	for x = 3, 6 do
		local f = NewFoundation(x, 0)
		SetCompareFunction(f, "Append", "UpSuitWrap")
	end
	for x = 0, 6 do
		NewReserve(x, 1, FAN_NONE)
	end
	for x = 0, 6 do
		local t = NewTableau(x, 2, FAN_DOWN, MOVE_ANY)
		SetCompareFunction(t, "Append", "DownAltColorWrap")
		SetCompareFunction(t, "Move", "DownAltColorWrap")
	end
end

function StartGame()
	local stock = Stock()
	for _, r in ipairs(Reserves()) do
		MoveCard(stock, r)
	end

	local fs = Foundations()
	local c = MoveCard(stock, fs[1])
	local ord = Ordinal(c)
	for _, f in ipairs(fs) do
		SetLabel(f, ordToChar[ord])
	end
	ord = ord - 1
	if ord == 0 then
		ord = 13
	end
	for _, t in ipairs(Tableaux()) do
		SetLabel(t, ordToChar[ord])
	end

	-- Klondike tableaux
	local dealDown = -1
	for _, t in ipairs(Tableaux()) do
		for _ = 0, dealDown do
			c = MoveCard(stock, t)
			FlipDown(c)
		end
		dealDown = dealDown + 1
		MoveCard(stock, t)
	end
end

function TailTapped(tail)
	if Category(Owner(First(tail))) == "Stock" then
		local stock = Stock()
		for _, r in ipairs(Reserves()) do
			MoveCard(stock, r)
		end
	else
		DefaultTailTapped(tail)
	end
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Agnes_(solitaire)"
end
