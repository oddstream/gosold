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
	end

end

function StartGame(moonHandle)
	print("Hello from lua StartGame")
end

-- function TailMoveError(moonHandle, tail)
-- 	print("Hello from lua TailMoveError")
-- 	return false, "Cannot move that tail"
-- end

print("Lua FreecellScript loaded")
