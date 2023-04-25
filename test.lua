
function BuildPiles(moonHandle)
	print("Hello from lua BuildPiles")
end

function StartGame(moonHandle)
	print("Hello from lua StartGame")
end

function TailMoveError(moonHandle, tail)
	print("Hello from lua TailMoveError")
	return false, "Cannot move that tail"
end

print("Lua test loaded")
