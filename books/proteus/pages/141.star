load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
This will be a bitter and bloody battle to the
death:

%s
""" % fight_render("BELENGHAST", 14, 22)

on_command = fight_build_command({
    "goto": "198",
})
