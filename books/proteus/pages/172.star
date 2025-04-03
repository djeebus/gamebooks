load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
The potion takes effect: you feel immensely
strong, and for the first two rounds of combat,
each blow you land will cause four Strength
points of damage to the GIANT. After that,
each blow will cause only 2 Strength points of
damage, as normal. Each blow to you will cost
you 2 Strength points throughout the battle.

%s
""" % fight_render("GIANT", 12, 12)

on_command = fight_build_command({
    "goto": "126",
    "item_id": "power-potion",
})
