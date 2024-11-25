load("../lib/stats.star", "stat_courage")

markdown = """
The tunnel goes South, then after a while
turns West. You have gone West only a short
distance when your foot catches on something
and you stumble. [Roll two dice](!next). If the score is
less than your current Courage score, turn to 88. 
If the score is the same or greater, turn to 14.
"""

def on_command(cmd, args):
    if cmd == "next":
        score = dice_roll(2, 6)
        if score < stat_courage():
            return "88"
        return "14"
