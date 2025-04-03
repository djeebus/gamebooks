load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
You fight the SCRAFE.

%s
""" % fight_render("SCRAFE", 10, 12)

on_command = fight_build_command({
    "goto": "152",
})
