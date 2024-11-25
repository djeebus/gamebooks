load("../lib/fight.star", "fight_render", "fight_build_command")
load("../lib/stats.star", "strength_add")

def once():
    strength_add(-2)


markdown = """
You do not move quickly enough, and the
lance gashes your side. Lose 2 Strength
points. The TROLL looks strong, and cunning,
 and he is quick with his lance. He makes
another sudden lunge, and you parry it with
your shield just in time.

You fight the TROLL.

%s
""" % fight_render("TROLL", 10, 10)

on_command = fight_build_command({
    "goto": "76",
})
