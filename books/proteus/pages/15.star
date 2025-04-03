load("../lib/fight.star", "fight_build_command", "fight_render")

markdown = """As you drink the potion, the GIANT lets out
a deafening shriek, smashes his fist into the
rock wall, then suddenly charges at you,
roaring and screaming. You barely have time
to realise that the Potion was not a wise choice
before you are locked in battle with the
enraged GIANT!

%s
""" % (fight_render("GIANT", 12, 18))

on_command = fight_build_command({"goto": "126"})
