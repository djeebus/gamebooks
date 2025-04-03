load("../lib/stats.star", "strength_add")

def once():
    score = dice_roll(1, 6)
    strength_add(-score)

markdown = """
As you pull the levers, there is a sudden, terrible
flash of pain. You drop to the ground. Roll one
dice: the score is the number of
Strength points you have lost. If you are still
alive, you may try again. Will you try:

- The left and middle levers? [Turn to 19](19)
- The left and right levers? [Turn to 140](140)
"""
