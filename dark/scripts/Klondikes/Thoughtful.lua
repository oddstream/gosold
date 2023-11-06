-- Thoughtful

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
	-- Klondike face up tableaux
	local deal = 0
	for _, t in ipairs(Tableaux()) do
		for _ = 0, deal do
			MoveCard(stock, t)
		end
		deal = deal + 1
	end
	SetRecycles(2)
	populateEmptyWasteFromStock()
end

function AfterMove()
	populateEmptyWasteFromStock()
end

function PileTapped(pile)
	if Category(pile) == "Stock" then
		recycleWasteToStock()
	end
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Klondike_(solitaire)"
end
