-- Usk

local layout = {
	{x=0, n=8},
	{x=1, n=8},
	{x=2, n=8},
	{x=3, n=7},
	{x=4, n=6},
	{x=5, n=5},
	{x=6, n=4},
	{x=7, n=3},
	{x=8, n=2},
	{x=9, n=1},
}

function BuildPiles()
	NewStock(0, 0)
	for x = 6, 9 do
		local f = NewFoundation(x, 0)
		SetCompareFunction(f, "Append", "UpSuit")
		SetLabel(f, "A")
	end
	for _, li in ipairs(layout) do
		local t = NewTableau(li.x, 1, FAN_DOWN, MOVE_ANY)
		SetCompareFunction(t, "Append", "DownAltColor")
		SetCompareFunction(t, "Move", "DownAltColor")
		SetLabel(t, "K")
	end
end

local function dealCards()
	local stock = Stock()
	local tx = Tableaux()
	for i, li in ipairs(layout) do
		local t = tx[i]
		for _ = 1, li.n do
			MoveCard(stock, t)
		end
	end
end

function StartGame()
	dealCards()
	SetRecycles(1)
end

function PileTapped(pile)
	if Category(pile) ~= "Stock" then
		return
	end
	if Recycles() == 0 then
		Toast("No more recycles")
		return
	end
	for _, t in ipairs(Tableaux()) do
		if Len(t) > 0 then
			MoveTail(First(t), Stock())
		end
	end
	Reverse(Stock())
	dealCards()
	SetRecycles(0)
end

function Wikipedia()
	return "https://politaire.com/help/usk"
end