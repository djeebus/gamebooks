local dice = require("gamebooks/dice")
local storage = require("gamebooks/storage")

local party_key = "--party-key--"
local cunning_key = "--cunning--"

local _M = {
    name = "Star smuggler",
    start_page = "e001",

    make_cunning_roll = function()
        local cunning_score = storage.get(cunning_key)
        if cunning_score == nil then
            error("cunning score is missing!")
        end
        local roll = dice.roll(1, 6)
        local success = roll <= cunning_score
        return success
    end,

    wrap_content = function(body)
        local content = [[
# Cool header content
]] .. body
        content = content .. [[
# Cool footer content
        ]]
        return content
    end,
}

function _M.on_start()
    -- set initial cunning score
    local cunning_score = dice.roll(1, 6)
    storage.set(cunning_key, cunning_score)
end

function _M.get_party()
    party = storage.get(party_key)
    if party == nil then
        party = {}
    end
    return party
end

return _M
