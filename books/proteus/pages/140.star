load("../lib/stats.star", "strength_add")

def once():
    score = dice_roll(1, 6)
    strength_add(-score)

markdown = """
There is a flash of blinding bright pain as
you pull the levers, and you drop, halfconscious,
to the ground. Roll one dice: the
score is the number of Strength points you
have lost. If you are still alive, you may try
again. Will you now try:

- The left and middle levers? [Turn to 19](19)
- The middle and right levers? [Turn to 107](107)
"""
