-- Selective Canfield

local ordToChar = {"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

local function populateEmptyWasteFromStock()
	local stock = Stock()
	local waste = Waste()
	if Len(waste) == 0 then
		MoveCard(stock, waste)
		MoveCard(stock, waste)
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
		-- When moving tableau piles, you must either move the whole pile or only the top card.
		local t = NewTableau(x, 1, FAN_DOWN, MOVE_ONE_OR_ALL)
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

	-- 20 cards in the reserve
	local r = Reserve()
	for _ = 1, 19 do
		local c = MoveCard(stock, r)
		FlipDown(c)
	end
	MoveCard(stock, r) -- last card face up

	-- one card in each tableau
	local tx = Tableaux()
	for _, t in ipairs(tx) do
		MoveCard(stock, t)
	end

	populateEmptyWasteFromStock()
	SetRecycles(32767)

	Toast("Move a Tableaux card to a Foundation")

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
			Toast("Move a Tableaux card to a Foundation")
		else
			for _, f in ipairs(fs) do
				SetLabel(f, ordToChar[ord])
			end
		end
	end

	-- Empty tableaux spaces are filled automatically from the reserve.
	for _, t in ipairs(Tableaux()) do
		if Len(t) == 0 then
			MoveCard(Reserve(), t)
		end
	end
	populateEmptyWasteFromStock()
end

-- default TailMoveError

-- default TailAppendError
	-- Once the reserve is empty, spaces in the tableau can be filled with a card from the Deck [Stock/Waste],
	-- but NOT from another tableau pile.
	-- pointless rule, since tableaux move rule is MOVE_ONE_OR_ALL

function TailTapped(tail)
	local card = First(tail)
	local pile = Owner(card)
	if Category(pile) == "Stock" then
		MoveCard(Stock(), Waste())
		MoveCard(Stock(), Waste())
		MoveCard(Stock(), Waste())
	else
		DefaultTailTapped(tail)
	end
end

function PileTapped(pile)
	if Category(pile) == "Stock" then
		recycleWasteToStock()
	end
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Canfield_(solitaire)"
end
