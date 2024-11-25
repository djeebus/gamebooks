load("../lib/stats.star", "stat_courage")

markdown = """
The silver key fits the lock and turns easily,
and you push open the door. You are in a small
stone room, and you quickly realise that even
the sparse furnishings – the table, chair, and
bed – are also of stone.

A snuffling, grunting sound from your left
makes you wheel round, and you are face to
face with a TROLL. The TROLL is a small,
squat creature, smaller than you, but as broad
across. He is bald, and has an ugly, fearsome
face.

He moves towards you, wielding a short
silver lance, which he suddenly thrusts at you.
[Roll 2 dice](!next).
"""

def on_command(cmd, args):
    if cmd == "next":
        score = dice_roll(2, 6)
        courage = stat_courage()
        if score <= courage:
            return "101"
        else:
            return "113"
