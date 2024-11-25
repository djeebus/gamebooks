load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
The STONEMAN closes with your lookalike. Throw for each of them: your look-alike
has your current Strength and Dexterity
scores. Only if the STONEMAN defeats your
look-alike will you then have to finish him off,
with your Strength and Dexterity scores as
they were before this battle.

%s
""" % fight_render("STONEMAN", 12, 14)

on_command = fight_build_command({
    "item_id": "duality-potion",
    "goto": "117",
})
