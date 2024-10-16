local _M = {
    title = "The Adventure Begins",
    duration = 0,
    markdown = [[
Due to bad luck and loan sharks, your financial situation is
getting very desperate. Your small merchant starship never
seems to have a full cargo hold, but operating costs are high.
Your cash is almost gone, and another payment is due to the
sector financiers. Maybe, just maybe, you can make ends meet
if you look for illegal goods, and begin to take chances.
Consult [r201b] to determine your starting attributes, and your
skills, for your new career as a Star Smuggler, then continue
reading here:
You have a sturdy and reliable [Antelope class starship](r212)
built to [tech level 1](r210) standards and outfitted with a
[Hopper class ship's boat](r214) and [starship guns](r216a), both
also tech level 1. The starship has six [Hypercharges](r212b) and
the boat has 15 [fuel units](r211). In addition, a [stasis unit](r212e)
is mounted in the pilot's compartment with 2 CU
capacity, to protect the occupants in case of disaster. You
personally own a [utility suit](r213) and a [sidearm](r216d), both
of tech level 1. Your only money is in your pocket: 1d6 times
100, plus 150 Secs (for example, a 1d6 roll of "3" means 3x100,
then +150, or 450 Secs, see r232a).
You have no crew or hirelings, no cargo, and no repair units.
However, you do have proper papers and are not "wanted" in
any system (r228). This is the first day of the week, so you have
10 days until your next starship payment is due (r203e).
You are currently at the sole planet in the Regari system
(r207a) of the Pavonis sector, at the spaceport (r205o). This
morning you decide to take up a life of smuggling. You check
over your starship guns and personal sidearm, and prepare to
find profit through any means. See r203 for the activities
available to you.
    ]],
}

function _M.on_enter()
    get_or_set_roll("initial_cache", "(1d6*100)+150")
end

return _M
