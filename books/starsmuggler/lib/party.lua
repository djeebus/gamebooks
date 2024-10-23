local storage = require("gamebooks/storage")

local party_key = "--party-key--"

return {
    get_party = function()
        party = storage.get(party_key)
        if party == nil then
            party = {}
        end
        return party
    end,
}
