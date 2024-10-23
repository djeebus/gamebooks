local storage = require("gamebooks/storage")
local dice = require("gamebooks/dice")

local cash_key = '--cash-key--'

return {
    init = function()
        initial_cash = dice.roll(1, 6) * 100 + 150
        storage.set(cash_key, initial_cash)
    end,
    get_cash = function()
        return storage.get(cash_key)
    end
}
