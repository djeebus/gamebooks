load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
You will lose only one Strength point for the
first three times that the STONEMAN scores
a blow on you in this battle. Thereafter, you will
lose the usual 2 Strength points when the
STONEMAN scores a blow against you.

%s

""" % fight_render("STONEMAN", 12, 14)

on_command = fight_build_command({
    "item_id": "elusiveness-potion",
    "goto": "117",
})
