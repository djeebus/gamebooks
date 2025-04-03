load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
You feel tremendously powerful. For the
first two rounds of this battle, each blow you
score against the SCRAFE will cost it 4
Strength points. After that, each blow will
cause 2 points of damage. Each successful
attack by the SCRAFE against you will cost
you the normal 2 points of damage.

%s
""" % fight_render("SCRAFE", 10, 12)

on_command = fight_build_command({
    "goto": "152",
    "item_id": "power-potion",
})
