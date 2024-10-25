load("../lib/character.star", "profession_get")
load("../lib/god.star", "god_get")

def _build():
    difficulty = 10
    if profession_get() == "Wayfarer":
        difficulty -= 3
    if god_get() == "Lacuna":
        difficulty -= 1

    return """
You explain what a tree-lover you are, and that your heart has
always been at home in forests. You portray yourself as the
greenest adventurer ever.

[Make a CHARISMA roll at Difficulty 10](!roll!charisma!%d!391!410).
    """ % difficulty

markdown = _build()
