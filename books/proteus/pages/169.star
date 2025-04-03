load("../lib/stats.star", "strength_add", "courage_add")

def once():
    strength_add(-3)
    courage_add(1)


markdown = """
You have gone some way, and all is quiet.
But then the passage starts to narrow, and you
have to make your way sideways, holding
your torch ahead of you. A sudden fierce pain
in your leg stops you in your tracks, and you
see by the fight of your torch that your leg is
caught in the jaws of powerful metal mantrap.

Dropping the torch, you force the jaws
apart. Lose 3 Strength points but gain 1
Courage point.

Ignoring the pain, you pick up your torch
and continue East, aware that you must be
dripping blood. There is a passage South, but
it seems dark and ominous, and so you
continue on your way. You realise that the
rocks themselves are giving a shadowy, greenish 
light, which allows you to see better ahead.
[Turn to 134](134).
"""
