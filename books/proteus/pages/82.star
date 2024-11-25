load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
The door is as well-built as it looks, and only
after repeated blows from your shoulder
charges does it begin to splinter. Eventually,
you are able to wrench it off its hinges, but it
costs you 2 Strength points.

Beyond the door, there is a small, square
room, and in the centre of the floor is a statue.
You go closer, and see that the statue is of a
large, broad-shouldered man with a low forehead and a menacing expression. But then
you notice, at the foot of the statue, a crystal
key. Carefully, you reach for the key, and as
you do so, the statue begins moving!

You had half-expected this, and now leap
back, your sword raised, as the STONEMAN
moves towards you. If you have any of the
following Potions, and wish to use them here,
turn to the appropriate section:

- Duality [Turn to 124](124)
- Elusiveness [Turn to 48](48)
- Fear [Turn to 102](102)

If you have none of these, or do not wish to
use them here, you must fight the STONEMAN alone.

%s
""" % fight_render("STONEMAN", 12, 14)

on_command = fight_build_command({
    "goto": 117
})
