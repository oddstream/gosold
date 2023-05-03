-- Whitehead

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
		SetCompareFunction(t, "Append", "DownColor")
		SetCompareFunction(t, "Move", "DownSuit")
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
	SetRecycles(0)
	populateEmptyWasteFromStock()
end

function AfterMove()
	populateEmptyWasteFromStock()
end

function Wikipedia()
	return "https://en.wikipedia.org/wiki/Klondike_(solitaire)"
end
