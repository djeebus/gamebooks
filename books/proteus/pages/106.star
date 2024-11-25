load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
Slightly to your surprise, the Potion seems to
be working! BELENGHAST roars, smashing
wildly with his axe; he seems blinded, and
you quickly close in battle. You get in three
good blows with your sword, but his armour
is good, and he is only slightly wounded.

The Potion begins to wear off, and he
becomes calmer. He turns to you with murder
in his eyes, the great axe whistling over his
head. You have weakened him, but now you
must fight to the death:

%s
""" % fight_render("BELENGHAST", 11, 17)

on_command = fight_build_command({
    "goto": "198",
})
