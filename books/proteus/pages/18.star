load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
The potion works – you are now very
difficult to wound. You will lose only 1
Strength point, instead of 2, for the first three
successful attacks against you. After that, the
potion wears off, and each successful attack
against you by the SCRAFE will cause the
usual 2 points of damage.

%s

If you win, turn to 152.
""" % (fight_render("SCRAFE", 10, 12))


on_command = fight_build_command({
    "item_id": "elusiveness-potion",
    "goto": "152",
})
