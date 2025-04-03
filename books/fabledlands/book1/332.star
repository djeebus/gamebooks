load("../lib/bank.star", "bank_withdraw")
load("../lib/cargo.star", "furs", "grain", "metals", "minerals", "spices", "textiles", "timber")
load("../lib/market.star", "market_render")
load("../lib/ship.star", ship_assert_ownership="assert_ownership")

items = {
    "furs": {
        "item": furs,
        "buy": 135,
        "sell": 130,
    },
    "grain": {
        "item": grain,
        "buy": 220,
        "sell": 210,
    },
    "metals": {
        "item": metals,
        "buy": 600,
        "sell": 570,
    },
    "minerals": {
        "item": minerals,
        "buy": 400,
        "sell": 310,
    },
    "spices": {
        "item": spices,
        "buy": 900,
        "sell": 820,
    },
    "textiles": {
        "item": textiles,
        "buy": 325,
        "sell": 285,
    },
    "timber": {
        "item": timber,
        "buy": 120,
        "sell": 100,
    },
}

markdown = """
The quayside at the Trading Post is a simple affair, without the
harbour facilities or the shipwrights of the larger ports. There are
no ships available to buy here, but if you already own a ship,
you can buy and sell cargo.

%s

Prices are for single Cargo Units. 

You can buy a one-way passage to [Yellowport](!passage) from here at a cost of 15 Shards.

If you own a ship, you can [set sail](!setsail). 

[Leave](195)
""" % (market_render("Cargo", items, "furs", "grain", "metals", "minerals", "spices", "textiles", "timber"))


def on_command(command):
    if command == "passage":
        bank_withdraw(15)
        return "74"

    if command == "setsail":
        ship_assert_ownership()
        return "155"

