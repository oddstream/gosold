-- FreecellScript

function BuildPiles(moonHandle)
	print("Hello from lua BuildPiles")
	NewStock(moonHandle, -5, -5)
	for x = 0, 3 do
		NewCell(moonHandle, x, 0)
	end
	for x = 4, 7 do
		local f = NewFoundation(moonHandle, x, 0)
		SetCompareFunction(moonHandle, f, "Append", "UpSuit")
		SetLabel(moonHandle, f, "A")
	end
	for x = 0, 7 do
		local t = NewTableau(moonHandle, x, 1, FAN_DOWN, MOVE_ONE_PLUS)
		SetCompareFunction(moonHandle, t, "Append", "DownAltColor")
		SetCompareFunction(moonHandle, t, "Move", "DownAltColor")
	end

end

function StartGame(moonHandle)
	print("Hello from lua StartGame")
	local stock = GetStock(moonHandle)
	-- 4 tabs [0 .. 3] with 7 cards
	-- 4 tabs [4 .. 7] with 6 cards
	for i = 0, 3 do
		local tab = GetTableau(moonHandle, i)
		for _ = 1, 7 do
			MoveCard(moonHandle, stock, tab)
		end
	end
	for i = 4, 7 do
		local tab = GetTableau(moonHandle, i)
		for _ = 1, 6 do
			MoveCard(moonHandle, stock, tab)
		end
	end
end

-- function TailMoveError(moonHandle, tail)
-- 	print("Hello from lua TailMoveError")
-- 	return false, "Cannot move that tail"
-- end

print("Lua FreecellScript loaded")
