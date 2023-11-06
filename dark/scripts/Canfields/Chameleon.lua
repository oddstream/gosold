-- Chameleon

local ordToChar = {"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

local function populateEmptyWasteFromStock()
	local stock = Stock()
	local waste = Waste()
	if Len(waste) == 0 then
		MoveCard(stock, waste)
	end
end

local function recycleWasteToStock()
	local recycles = Recycles()
	if recycles > 0 then
		local stock = Stock()
		local waste = Waste()
		while Len(waste) > 0 do
			MoveCard(waste, stock)
		end
		SetRecycles(recycles - 1)
	else
		Toast("No more recycles")
	end
end

function BuildPiles()

	NewStock(5, 0)

	NewWaste(5, 1, FAN_DOWN3)

	NewReserve(0, 1, FAN_NONE)

	for x = 0, 3 do
		local f = NewFoundation(x, 0)
		SetCompareFunction(f, "Append", "UpSuitWrap")
	end

	for x = 1, 3 do
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ANY)
		SetCompareFunction(t, "Append", "DownWrap")
		SetCompareFunction(t, "Move", "DownWrap")
	end

end

function StartGame()

	local stock = Stock()

	local r = Reserve()
	for _ = 1, 12 do
		MoveCard(stock, r)
	end

	local fs = Foundations()
	local card = MoveCard(stock, fs[1])
	local ord = Ordinal(card)
	for _, f in ipairs(fs) do
		SetLabel(f, ordToChar[ord])
	end

	for _, t in ipairs(Tableaux()) do
		MoveCard(stock, t)
	end

	SetRecycles(0)

	populateEmptyWasteFromStock()
end

function AfterMove()
	-- "fill each [tableau] space at once with the top card of the reserve,
	-- after the reserve is exhausted, fill spaces from the waste pile,
	-- but at this time a space may be kept open for as long as desired"
	for _, t in ipairs(Tableaux()) do
		if Len(t) == 0 then
			MoveCard(Reserve(), t)
		end
	end
	populateEmptyWasteFromStock()
end

-- default TailMoveError
-- default TailAppendError
-- default TailTapped

function PileTapped(pile)
	if Category(pile) == "Stock" then
		recycleWasteToStock()
	end
end

function CardColors()
	return 2
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Canfield_(solitaire)"
end
