load("../lib/stats.star", "stat_courage")


markdown = """
You run for the archway as the ZOMBIE
lurches towards you. [Roll 2 dice](!next). If the score is
less than your current Courage score, 
[turn to 47](47). If it is the same or 
greater, [turn to 171](171).
"""

def on_command(cmd, args):
    if cmd == "next":
        score = roll_dice()
        if score < stat_courage():
            return "47"
        else:
            return "171"
