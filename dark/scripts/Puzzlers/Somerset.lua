-- Somerset

local layout = {
	{x=0, n=1},
	{x=1, n=2},
	{x=2, n=3},
	{x=3, n=4},
	{x=4, n=5},
	{x=5, n=6},
	{x=6, n=7},
	{x=7, n=8},
	{x=8, n=8},
	{x=9, n=8},
}

function BuildPiles()
	NewStock(-5, -5)
	for x = 6, 9 do
		local f = NewFoundation(x, 0)
		SetCompareFunction(f, "Append", "UpSuit")
		SetLabel(f, "A")
	end
	for _, li in ipairs(layout) do
		local t = NewTableau(li.x, 1, FAN_DOWN, MOVE_ONE_PLUS)
		SetCompareFunction(t, "Append", "DownAltColor")
		SetCompareFunction(t, "Move", "DownAltColor")
	end
end

function StartGame()
	local stock = Stock()
	local tx = Tableaux()
	for i, li in ipairs(layout) do
		local t = tx[i]
		for _ = 1, li.n do
			MoveCard(stock, t)
		end
	end
	SetRecycles(0)
end

function Wikipedia()
	return "https://politaire.com/help/somerset"
end