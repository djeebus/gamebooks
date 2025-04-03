load("../lib/stats.star", "courage_add")


def once():
    courage_add(-2)


markdown = """
You walk Eastwards for a while, until your
way is blocked by a door. It is a solid wooden
door, covered in cobwebs, and does not look
very inviting; however, you push tentatively at
it, and it opens. The light of your torch barely
penetrates the thick, suffocating blackness
inside, but you can just make out an archway
in the wall opposite.

As you move towards it, you stumble on
something on the floor. You recoil – it looks
like a corpse that has been there for some
time, the flesh rotting away from the bones
and maggots crawling out of its mouth. As
you instinctively move away, its eyes suddenly 
flick open, and it moves towards you. The
door behind you slams shut. Lose 2 Courage
points as you back away in terror.

You cut and slash at it with your sword, but
it has no effect. The ZOMBIE is not vulnerable
to ordinary weapons, and no Potion of Magic
is any use against this creature. Do you have a
silver lance? If so, [turn to 29](!must-have-item!silver-lance!29). Otherwise, [turn
to 165](165).
"""
