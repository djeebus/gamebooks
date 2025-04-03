load("../lib/fight.star", "fight_multiple_render", "fight_build_command")

markdown = """
As you walk East, you hear a faint fluttering
in the darkness ahead of you. You press on,
your sword and shield at the ready, until you
realise that dozens of flying creatures are
swooping and diving around you. By the light
from your torch, you see that each flying
creature has the body of a huge bat, with
razor-sharp talons, but the face of a Ghoul.
You swing your sword, and most of them
retreat into the darkness, but three continue
the attack. Fight each in turn. If you have a
potion of Invincibility, Elusiveness, or Power,
you may use it here if you wish – but only
against the first DEATHBAT.

%s

""" % fight_multiple_render([
    {"name": "FIRST DEATHBAT", "dexterity": 11, "strength": 8},
    {"name": "SECOND DEATHBAT", "dexterity": 10, "strength": 6},
    {"name": "THIRD DEATHBAT", "dexterity": 8, "strength": 8},
])

on_command = fight_build_command({
    "goto": "6",
})
