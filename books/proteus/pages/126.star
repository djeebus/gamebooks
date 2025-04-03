load("../lib/inventory.star", "inventory_add")

def once():
    inventory_add("silver-key")

markdown = """
At the dead GIANT’s waist, you see a small
leather pouch. Inside is a silver key. Feeling
that this may be of use to you, you put it in
your pocket and continue North. Shortly, you
realise you are at a junction, and may either
continue North, or take a new passage that
leads off to your right. Will you:

- Take the new passage East? [Turn to 3](3)
- Continue North? [Turn to 182](182)"""
