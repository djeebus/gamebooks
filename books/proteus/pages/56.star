load("../lib/stats.star", "strength_add")


def once():
    strength_add(-3)


markdown = """
You have been going North for a little while,
when you suddenly cry out in pain. You twist
your head to see a poison dart in your neck.
Lose 3 Strength points. If you are still alive,
you wrench it out and press on warily North.
You come to a passage leading off to your left.

Will you:

- Take this way West? [Turn to 79](79)
- Continue North? [Turn to 145](145)
"""
