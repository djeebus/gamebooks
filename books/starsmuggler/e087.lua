local table = require("table")
local book = require("game")
local dice = require("gamebooks/dice")
local storage = require("gamebooks/storage")

local title = "Prisoner Pleads for Escape"

party = book.get_party()


if #party > 0 then
    return {
        title = title,
        markdown = "If you are with one or more people as a party, ignore this event.",
    }
end
if storage.get_page("step-1") == nil then
    text = [[
    a prisoner approaches the bars and wall to
    address you. He says that if you get him out and take him to
    any Rural, Rough, or Slum area on Byzantium or Imperia he can
    pay you 8,000 S.
    ]]

    storage.set_page("step-1", true)

    -- is the prisoner telling the truth? store for later
    local truth_roll = dice.roll(1, 6)
    local is_truth = (truth_roll ~= 6)
    storage.set("e087-prisoner-truth", is_truth)

    -- do you know if he's telling the truth?
    local success = book.make_cunning_roll()  -- r202 for implementation
    if success then
        if is_truth then
            text = text .. " You can tell that he's telling the truth."
        else
            text = text .. " You can tell that he's lying; you'll get nothing for your pains."
        end
    end

    text = text .. [[ To get him out takes the rest of the day, and requires that you
[disable or kill a guard](r301) with E7,M 1d6,H5 with a sidearm (see
r210 for tech level). You then escape,
and will leave this area and enter another as part of the activity
(with an appropriate entry encounter).
The former prisoner is "wanted" in this star system, and is E
1d6, M 1d6, H 1d6+1 (maximum of 6). He is a willing member
of your party until you next reach Byzantium or Imperia.]]

    return {
        title = title,
        markdown = text,
    }
end

if storage.get_page("step-2") == nil then
    do_you_accept = storage.get_page("accept")
    if do_you_accept ~= nil and do_you_accept == "yes" then
        return {
            title = title,
            markdown = [[
                do something
            ]]
        }
    end

    return {
        title = title,
        markdown = [[
            do something else
        ]]
    }
end

local page = {
    title = "Prisoner Pleads for Escape",
    metadata = {
        duration = 1,
    },
    quest = {
        {
            continue_if = function()
                party = book.get_party()
                return table.getn(party) == 0
            end,
            text = "If you are with one or more people as a party, ignore this event.",
        },
        {
            text = [[
                If you are alone, a prisoner approaches the bars and wall to
address you. He says that if you get him out and take him to
any Rural, Rough, or Slum area on Byzantium or Imperia he can
pay you 8,000 S. He is telling the truth on a result of 1-5 with a
1d6 roll. If he is lying, you'll get nothing for your pains.
            ]],
        },
        {
            action = {
                type = "roll_die",
                die_count = 1,
                die_size = 6,
                key = lie_die_storage_key,
            },
        },
        {
            condition_if = function()
                die = storage.get(lie_die_storage_key)
                return die == 6
            end,
            text = "If he is lying, you'll get nothing for your pains.",
        },
        {}
    },
    text = [[
 You can determine that now if you make a successful
[Cunning roll](r202), otherwise you don't find out the truth until
you deliver him. If he is lying, you'll get nothing for your pains.
To get him out takes the rest of the day, and requires that you
disable or kill a guard with E7,M 1d6,H5 with a sidearm (see
r210 for tech level). See r301 for the attack. You then escape,
and will leave this area and enter another as part of the activity
(with an appropriate entry encounter).
The former prisoner is "wanted" in this star system, and is E
1d6, M 1d6, H 1d6+1 (maximum of 6). He is a willing member
of your party until you next reach Byzantium or Imperia.
]]
}

return page
