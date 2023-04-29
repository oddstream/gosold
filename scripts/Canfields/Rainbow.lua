-- Rainbow

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

	NewStock(0, 0)

	NewWaste(1, 0, FAN_RIGHT3)

	NewReserve(0, 1, FAN_DOWN)

	for x = 3, 6 do
		local f = NewFoundation(x, 0)
		SetCompareFunction(f, "Append", "UpSuitWrap")
	end

	for x = 3, 6 do
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ONE_OR_ALL)
		SetCompareFunction(t, "Append", "DownWrap")
		SetCompareFunction(t, "Move", "DownWrap")
	end

end

function StartGame()

	local stock = Stock()

	local r = Reserve()
	for _ = 1, 12 do
		local c = MoveCard(stock, r)
		FlipDown(c)
	end
	MoveCard(stock, r)	-- final card face up

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
	for _, t in ipairs(Tableaux()) do
		if Len(t) == 0 then
			MoveCard(Reserve(), t)
		end
	end
	populateEmptyWasteFromStock()
end

function TailAppendError(dst, tail)
	if Len(dst) == 0 then
		if Category(dst) == "Tableau" then
			local rescards = 0
			for _, r in ipairs(Reserves()) do
				rescards = rescards + Len(r)
			end
			if rescards > 0 then
				local card = First(tail)
				local pile = Owner(card)
				if Category(pile) ~= "Reserve" then
					return false, "An empty Tableau must be filled from a Reserve"
				end
			end
		end
		return CompareEmpty(dst, tail)
	end
	return CompareAppend(dst, tail)
end

function PileTapped(pile)
	if Category(pile) == "Stock" then
		recycleWasteToStock()
	end
end

function CardColors()
	return 4
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Canfield_(solitaire)"
end
