load("../lib/items.star", "armor1", "armor2", "armor3", "armor4", "armor5", "armor6")
load("../lib/market.star", "market_render")

items = {
    "armour1": {
        "item": armor1,
        "buy": 50,
        "sell": 45,
    },
    "armour2": {
        "item": armor2,
        "buy": 100,
        "sell": 90,
    },
    "armour3": {
        "item": armor3,
    }
}

the_market = "\n\n".join([
    (market_render("Armour", items, "armour1", "armour2", "armour3", "armour4", "armour5", "armour6"))
])

markdown = """
The market is large and busy. At the corner of Brimstone Plaza,
gigantic braziers burn sweet-smelling incense in an attempt to
overpower the rotten-egg smell that permeates the whole city.
There are many stalls and goods to chose from. You may buy
any of the items listed, as long as you have the money and the
space to carry it. You may also sell any items you own that are
listed below, for the price stated – if you do, don’t forget to
cross them off your Adventure Sheet.

Items with no purchase price are not available locally.

%s

One trader is offering a treasure map for sale at 200 Shards,
and will buy any old treasure map for 150 Shards. If you buy
the map, note this paragraph number (30) for reference and turn
to 200.

If you wish to buy cargo for a ship, you need to visit the
warehouses at the harbourmaster. When you are ready to return
to the city centre, [turn to 10](10).
""" % (the_market)
