-- Australian

local function populateEmptyWasteFromStock()
	local stock = Stock()
	local waste = Waste()
	if Len(waste) == 0 then
		MoveCard(stock, waste)
	end
end

function BuildPiles()

	NewStock(0, 0)

	NewWaste(1, 0, FAN_RIGHT3)

	for x = 4, 7 do
		local f = NewFoundation(x, 0)
		SetCompareFunction(f, "Append", "UpSuit")
		SetLabel(f, "A")
	end

	for x = 0, 7 do
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ANY)
		SetCompareFunction(t, "Append", "DownSuit")
		SetLabel(t, "K")
	end

end

function StartGame()

	local stock = Stock()

	local tx = Tableaux()
	for _, t in ipairs(tx) do
		for _ = 1,4 do
			MoveCard(stock, t)
		end
	end

	populateEmptyWasteFromStock()
	SetRecycles(0)

end

function AfterMove()
	populateEmptyWasteFromStock()
end

function TailMoveError(tail)
	-- override default (which tests if tail is conformant)
	return true, ""
end

-- default TailAppendError
-- default TailTapped
-- default PileTapped

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Australian_Patience"
end

function CardColors()
	return 4
end