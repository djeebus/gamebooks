load("../lib/stats.star", "strength_add")

def once():
    strength_add(-2)

markdown = """
The passage continues North for a long
way. After a while it becomes narrower, and
begins sloping upwards, and it starts to become 
much hotter. You wipe the sweat out of
your eyes and continue, but it feels as though
the whole tunnel is shaking around you, while
the heat increases. After a few seconds, you
realise that this is no illusion – the ground and
the walls are indeed trembling and shaking.
The shaking suddenly becomes violent, and
you are hurled against the wall of the tunnel.
Lose 2 Strength points.
As the shaking continues, you can see that
ahead of you, the ground has opened up –
there is a huge chasm. The shaking continues;
if you have a Potion of Flying, and wish to use
it now, [turn to 11](!must-have-item!flying-potion!11). 
Otherwise, you will have to
run back South, and take the passage East.
[Turn to 5](5).
"""
