local dice = require("gamebooks/dice")
local storage = require("gamebooks/storage")

local cunning_key = "--cunning--"

return {
    init = function()
        local cunning_score = dice.roll(1, 6)
        storage.set(cunning_key, cunning_score)
    end,

    make_cunning_roll = function()
        local cunning_score = storage.get(cunning_key)
        if cunning_score == nil then
            error("cunning score is missing!")
        end
        local roll = dice.roll(1, 6)
        local success = roll <= cunning_score
        return success
    end,
}
