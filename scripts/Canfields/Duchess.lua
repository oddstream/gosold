-- Duchess

local ordToChar = {"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

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

	NewStock(1, 1)

	for _, x in ipairs({0,2,4,6}) do
		NewReserve(x, 0, FAN_RIGHT)
	end

	NewWaste(1, 2, FAN_DOWN3)

	for x = 3, 6 do
		local f = NewFoundation(x, 1)
		SetCompareFunction(f, "Append", "UpSuitWrap")
	end

	for x = 3, 6 do
		local t = NewTableau(x, 2, FAN_DOWN, MOVE_ANY)
		SetCompareFunction(t, "Append", "DownAltColorWrap")
		SetCompareFunction(t, "Move", "DownAltColorWrap")
	end

end

function StartGame()

	local stock = Stock()

	local fs = Foundations()
	for _, f in ipairs(fs) do
		SetLabel(f, "")
	end

	local rs = Reserves()
	for _, r in ipairs(rs) do
		MoveCard(stock, r)
		MoveCard(stock, r)
		MoveCard(stock, r)
	end

	local ts = Tableaux()
	for _, t in ipairs(ts) do
		MoveCard(stock, t)
	end

	SetRecycles(1)
	Toast("Move a Reserve card to a Foundation")

end

function AfterMove()
	local fs = Foundations()
	if Label(fs[1]) == "" then
		local ord = 0
		for _, f in ipairs(fs) do
			if Len(f) > 0 then
				local card = Peek(f)
				ord = Ordinal(card)
				break
			end
		end
		if ord == 0 then
			Toast("Move a Reserve card to a Foundation")
		else
			for _, f in ipairs(fs) do
				SetLabel(f, ordToChar[ord])
			end
		end
	end
end

function TailAppendError(dst, tail)
	if Len(dst) == 0 then
		if Category(dst) == "Foundation" then
			if Label(dst) == "" then
				local card = First(tail)
				local pile = Owner(card)
				if Category(pile) ~= "Reserve" then
					return false, "The first Foundation card must come from a Reserve"
				end

			end
		elseif Category(dst) == "Tableau" then
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

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Duchess_(solitaire)"
end
