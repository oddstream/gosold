-- DuchessScript

function BuildPiles(moonHandle)

	NewStock(moonHandle, 1, 1)

	for _, x in ipairs({0,2,4,6}) do
		NewReserve(moonHandle, x, 0, FAN_RIGHT)
	end

	NewWaste(moonHandle, 1, 2, FAN_DOWN3)

	for x = 3, 6 do
		local f = NewFoundation(moonHandle, x, 1)
		SetCompareFunction(moonHandle, f, "Append", "UpSuitWrap")
	end

	for x = 3, 6 do
		local t = NewTableau(moonHandle, x, 2, FAN_DOWN, MOVE_ANY)
		SetCompareFunction(moonHandle, t, "Append", "DownAltColorWrap")
		SetCompareFunction(moonHandle, t, "Move", "DownAltColorWrap")
	end

end

function StartGame(moonHandle)

	local stock = GetStock(moonHandle)

	local fs = GetFoundations(moonHandle)
	for _, f in ipairs(fs) do
		SetLabel(moonHandle, f, "")
	end

	local rs = GetReserves(moonHandle)
	for _, r in ipairs(rs) do
		MoveCard(moonHandle, stock, r)
		MoveCard(moonHandle, stock, r)
		MoveCard(moonHandle, stock, r)
	end

	local ts = GetTableaux(moonHandle)
	for _, t in ipairs(ts) do
		MoveCard(moonHandle, stock, t)
	end

	SetRecycles(moonHandle, 1)
	Toast(moonHandle, "Move a Reserve card to a Foundation")

end

function AfterMove(moonHandle)
	local fs = GetFoundations(moonHandle)
	if Label(moonHandle, fs[1]) == "" then
		local ord = 0
		for _, f in ipairs(fs) do
			if Len(moonHandle, f) > 0 then
				local card = Peek(moonHandle, f)
				ord = Ordinal(moonHandle, card)
				break
			end
		end
		if ord == 0 then
			Toast(moonHandle, "Move a Reserve card to a Foundation")
		else
			local U = {"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
			for _, f in ipairs(fs) do
				SetLabel(moonHandle, f, U[ord])
			end
		end
	end
end

function TailAppendError(moonHandle, dst, tail)
	if Len(moonHandle, dst) == 0 then
		if Category(moonHandle, dst) == "Foundation" then
			if Label(moonHandle, dst) == "" then
				local card = First(moonHandle, tail)
				local pile = Owner(moonHandle, card)
				if Category(moonHandle, pile) ~= "Reserve" then
					return false, "The first Foundation card must come from a Reserve"
				end

			end
		elseif Category(moonHandle, dst) == "Tableau" then
			local rescards = 0
			for _, r in ipairs(GetReserves(moonHandle)) do
				rescards = rescards + Len(moonHandle, r)
			end
			if rescards > 0 then
				local card = First(moonHandle, tail)
				local pile = Owner(moonHandle, card)
				if Category(moonHandle, pile) ~= "Reserve" then
					return false, "An empty Tableau must be filled from a Reserve"
				end
			end
		end
		return CompareEmpty(moonHandle, dst, tail)
	end
	return CompareAppend(moonHandle, dst, tail)
end

function Wikipedia(moonHandle)
	return "https://en.wikipedia.org/wiki/Duchess_(solitaire)"
end
