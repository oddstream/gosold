-- Storehouse

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
		SetCompareFunction(t, "Append", "DownSuitWrap")
		SetCompareFunction(t, "Move", "DownSuitWrap")
	end

end

function StartGame()

	local stock = Stock()
	local founds = Foundations()

	local c
	c = Extract(stock, 0, 2, CLUB)
	Push(founds[1], c)
	c = Extract(stock, 0, 2, DIAMOND)
	Push(founds[2], c)
	c = Extract(stock, 0, 2, HEART)
	Push(founds[3], c)
	c = Extract(stock, 0, 2, SPADE)
	Push(founds[4], c)

	for _, f in ipairs(founds) do
		SetLabel(f, "2")	-- cosmetic; foundations will never be empty
	end

	local r = Reserve()
	for _ = 1, 12 do
		c = MoveCard(stock, r)
		FlipDown(c)
	end
	MoveCard(stock, r)	-- final card face up

	for _, t in ipairs(Tableaux()) do
		MoveCard(stock, t)
	end

	SetRecycles(2)

end

function AfterMove()
	for _, t in ipairs(Tableaux()) do
		if Len(t) == 0 then
			MoveCard(Reserve(), t)
		end
	end
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
		local recycles = Recycles()
		if recycles > 0 then
			local stock = Stock()
			local waste = Waste()
			while Len(waste) > 0 do
				MoveCard(waste, stock)
			end
			SetRecycles( recycles - 1)
		else
			Toast("No more recycles")
		end
	end
end

function CardColors()
	return 4
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Canfield_(solitaire)"
end
