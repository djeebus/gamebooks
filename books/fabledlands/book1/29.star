markdown = """
Your ship is sailing in the coastal waters beside Yellowport.
There are a number of other ships, mostly merchantmen, but
there are also a few warships of the Sokaran Imperial Navy.
‘At least we won’t be plagued by pirates with the navy
around,’ says the first mate.
[Roll two dice](!next)
Score 2-4 Storm turn to 613
Score 5-9 An uneventful voyage turn to 439
Score 10-12 Sokaran war galley turn to 165
"""

def on_command(command):
    if command == "next":
        result = dice_roll(2, 6)
        if result <= 4:
            return "613"
        elif result >= 10:
            return "165"
        else:
            return "439"
