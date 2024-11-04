load("../lib/items.star",
     "armor1", "armor2", "armor3", "armor4", "armor5", "armor6",
     "weapon0", "weapon1", "weapon2", "weapon3",
     "magical1", "magical2", "magical3",
     "lockpicks", "holysymbol", "compass", "parchment", "rope", "lantern")
load("../lib/market.star", "render", "on_command")

items = {
    "armor1": {
        "items": armor1,
        "buy": 50,
        "sell": 45,
    },
    "armor2": {
        "items": armor2,
        "buy": 100,
        "sell": 90,
    },
    "armor3": {
        "items": armor3,
        "buy": 200,
        "sell": 180,
    },
    "armor4": {
        "items": armor4,
        "sell": 360,
    },
    "armor5": {
        "items": armor5,
        "sell": 720,
    },
    "armor6": {
        "items": armor6,
        "sell": 1440,
    },
}

_market = "\n\n".join([
    render("Armour", items, "armor1", "armor2", "armor3", "armor4", "armor5", "armor6")
])

markdown = """
Wishport market is located in a warren of alleys adjacent to the
waterfront. Each alley is dedicated to a different commodity. In
one you can buy rope, in another lanterns, and so on. You hear
the echo of hammer on anvil from the weaponsmith’s forge just
ahead.

%s

Items which have no purchase price listed are not available locally.

The market is where you can buy and sell possessions to
carry on your person. To buy trade goods that can be carried on
board ship, you need to visit the warehouses at the waterfront.

[I'm done](217)
""" % _market
