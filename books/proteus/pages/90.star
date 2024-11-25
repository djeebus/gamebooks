load("../lib/stats.star", "strength_add")


def once():
    strength_add(-1)


markdown = """
You cannot see very far ahead by the light of
your torch, and have to move along the
passageway using your hands against the
walls to guide you. Your left hand suddenly
encounters nothing at all, and you fall heavily,
cutting your head. Lose 1 Strength point. By
the light of your torch you can make out a new
way North. Will you now:

- Continue East? [Turn to 120](120)
- Try the way North? [Turn to 45](45)
"""
