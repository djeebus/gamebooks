load("../lib/bank.star", "bank_balance", "bank_assert_min_balance")
load("../lib/inventory.star", "inventory_add")
load("../lib/items.star",
     "armor1", "armor2", "armor3",
     "weapon0", "weapon1",
     "magical1",
     "mandolin", "lockpicks", "holysymbol", "compass", "rope", "lantern",
     "climbinggear", "bagofpearls", "ratpoison", "pickaxe", "silvernugget"
     )
load("../lib/market.star", "market_render")

items = {
    "armor1": {
        "item": armor1,
        "buy": 50,
        "sell": 45,
    },
    "armor2": {
        "item": armor2,
        "sell": 90,
    },
    "armor3": {
        "item": armor3,
        "sell": 180,
    },
    "weapon0": {
        "item": weapon0,
        "buy": 50,
        "sell": 40,
    },
    "weapon1": {
        "item": weapon1,
        "sell": 200,
    },
    "magical1": {
        "item": magical1,
        "sell": 400,
    },
    "mandolin": {
        "item": mandolin,
        "buy": 300,
        "sell": 270,
    },
    "lockpicks": {
        "item": lockpicks,
        "sell": 270,
    },
    "holysymbol": {
        "item": holysymbol,
        "sell": 100,
    },
    "compass": {
        "item": compass,
        "buy": 500,
        "sell": 450,
    },
    "rope": {
        "item": rope,
        "buy": 50,
        "sell": 45,
    },
    "lantern": {
        "item": lantern,
        "buy": 100,
        "sell": 90,
    },
    "climbinggear": {
        "item": climbinggear,
        "sell": 90,
    },
    "bagofpearls": {
        "item": bagofpearls,
        "sell": 100,
    },
    "ratpoison": {
        "item": ratpoison,
        "sell": 50,
    },
    "pickaxe": {
        "item": pickaxe,
        "sell": 90,
    },
    "silvernugget": {
        "item": silvernugget,
        "sell": 150,
    },
}

_market = [
    ["Armour", items, "armor1", "armor2", "armor3"],
    ["Weapons (sword, axe, etc)", items, "weapon0", "weapon1"],
    ["Magical Equipment", items, "magical1"],
    [
        "Other items", items, "mandolin", "lockpicks", "holysymbol", "compass", "rope", "lantern", "climbinggear",
        "bagofpearls", "ratpoison", "pickaxe", "silvernugget",
    ],
]

markdown = """
The Trading Post has a small market place in the village square,
where half a dozen stalls sell a few goods. Items with no
purchase price are not available locally.

%s

[Return](195)
""" % (
    "\n\n".join([market_render(*args) for args in _market])
)


def on_command(command):
    if on_command(command):
        return

    fail("not implemented: %s" % command)
