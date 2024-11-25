load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
Your double suddenly appears, and moves in
for combat. But BELENGHAST merely
sneers, utters a few words, and then slices it in
two with the great battleaxe.

“Your magic is too weak, Stranger,” he says.
You close in final combat:

%s

Dexterity Strength
BELENGHAST 14 22
If you defeat him, turn to 198.
""" % fight_render("BELENGHAST", 14, 22)

on_command = fight_build_command({"goto": 198})
