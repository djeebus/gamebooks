load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
The passage continues West for a short
time, and then turns North. From the darkness 
ahead of you, you can hear the sound of
heavy footsteps, and you stop, your sword at
the ready.

Soon, into the flickering light, lumbers the
most enormous creature you have ever seen.
Man-like, twelve feet tall, and dressed in
skins, he carries a huge bone in his hand.
Seeing you, he stops, growls, and raises the
bone. You will have to fight the GIANT, but
you may use magic to help if you wish. If you
have any of these potions, you may take one
now:

- Power [Turn to 172](!consume-item!power-potion!172)
- Transparency [Turn to 191](!consume-item!transparency-potion!191)
- Madness [Turn to 15](!consume-item!madness-potion!15)

Otherwise, you will have to fight the GIANT unaided.

%s
""" % fight_render("GIANT", 12, 12)

on_command = fight_build_command({
    "goto": "126",
})
