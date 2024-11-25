load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
The STONEMAN cowers away, before returning to battle. But he is much less fearsome
now. You fight:

%s
""" % fight_render("STONEMAN", 9, 14)

on_command = fight_build_command({
    "goto": "117",
})
