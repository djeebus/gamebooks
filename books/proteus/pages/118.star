load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
The potion takes effect immediately and the
SCRAFE closes in combat with your double.
Throw for each of them: your double has your
current Dexterity and Strength scores. Only if
the SCRAFE defeats your double will you
have to fight it: in that case, your Dexterity and
Strength scores will be at the level they were
before this battle.

%s
""" % fight_render("SCRAFE", 10, 12)

on_command = fight_build_command({
    "item_id": "duality-potion",
    "goto": "152",
})
