load("../lib/stats.star", "courage_add")


def once():
    courage_add(-2)


markdown = """
You continue West for a while, until you
notice that you are passing under an archway
of stone. The darkness feels heavier, even
more oppressing, but you can just make out a
heavy wooden door in the wall opposite. You
move cautiously round the walls of the small
room you find yourself in, and try the door,
but it is firmly shut.

You move back towards the archway, but as
you do, your foot catches on something. You
stagger back – you have stepped on what
looks like a corpse that has been there for
some time: the flesh is rotting away from the
bones, and maggots are crawling out of its
mouth.

As you instinctively move away, its eyes
suddenly flick open, and it moves towards
you. You back away in terror – lose 2 Courage
points.

The yellow eyes gleam as the ZOMBIE
reaches out for you. You swing your sword,
but to no avail. The ZOMBIE is not vulnerable
to ordinary weapons, nor is any Potion of
Magic any use against this terrible creature.
Do you have a silver lance? If so, 
[turn to 29](!must-have-item!silver-lance!29). If
not, [turn to 165](165).
"""
