load("../lib/stamina.star", "stamina_add")

markdown = """
Soon you realize you are completely lost in this strange, magical
forest. You wander around for days, barely able to find enough
food and water. Lose 4 Stamina points. If you still live, you
eventually stagger out of the forest to the coast.

[Next page](!next)
"""

def on_command(command):
    if command == 'next':
        stamina_add(-4)
        return "128"
