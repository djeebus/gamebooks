load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
You leap aside, dodging the lance-thrust,
and close in battle:

%s
""" % fight_render("TROLL", 10, 10)

on_command = fight_build_command({
    "goto": "76",
})
