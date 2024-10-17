local storage = require("gamebooks/storage")
local dice = require("gamebooks/dice")

local initial_cash = storage.get("initial_cash")
if initial_cash == nil then
    initial_cash = dice.roll(1, 6) * 100 + 150
    storage.set("initial_cash", initial_cash)
    storage.set("cash", initial_cash)
end

local _M = {
    title = "The Adventure Begins",
    duration = 0,
    markdown = string.format([[
Due to bad luck and loan sharks, your financial situation is
getting very desperate. Your small merchant starship never
seems to have a full cargo hold, but operating costs are high.
Your cash is almost gone, and another payment is due to the
sector financiers. Maybe, just maybe, you can make ends meet
if you look for illegal goods, and begin to take chances.

Determine your [starting attributes and your
skills](r201b) for your new career as a Star Smuggler, then continue
reading here:

You have a sturdy and reliable [Antelope class starship](r212)
built to [tech level 1](r210) standards and outfitted with a
[Hopper class ship's boat](r214) and [starship guns](r216a), both
also tech level 1. The starship has six [Hypercharges](r212b) and
the boat has 15 [fuel units](r211). In addition, a [stasis unit](r212e)
is mounted in the pilot's compartment with 2 CU
capacity, to protect the occupants in case of disaster. You
personally own a [utility suit](r213) and a [sidearm](r216d), both
of tech level 1. Your only money is the %s Sees in your pocket.
You have no crew or hirelings, no cargo, and no repair units.
However, you do have proper papers and are not [wanted](r228) in
any system. This is the first day of the week, so you have
10 days until your next [starship payment is due](r203e).
You are currently at the sole planet in the [Regari system](r207a)
of the Pavonis sector, at [the spaceport](r205o). This
morning you decide to take up a life of smuggling. You check
over your starship guns and personal sidearm, and prepare to
find profit through any means. [Available activities](r203).
    ]], initial_cash)
}

return _M
