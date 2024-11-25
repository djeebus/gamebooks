load("../lib/stats.star", "stat_dexterity")

markdown = """
You realise that the tunnel is sloping downwards. Worse, the floor is covered with some
slimy substance, and you feel yourself starting
to lose your footing. [Roll two dice.](!next) If the score
is less than your Dexterity score, turn to 27. If it
is the same or greater, turn to 81.
"""

def on_command(cmd, args):
    if cmd == "next":
        score = dice_roll(2, 6)
        dexterity = stat_dexterity()
        if score < dexterity:
            return "27"
        return "81"
