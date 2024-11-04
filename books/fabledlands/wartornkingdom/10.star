load("../lib/codewords.star", "codeword_assert", "codeword_has")
load("../lib/house.star", "house_purchase")
load("../lib/inventory.star", "inventory_assert")

_track_visits_key = "track-visits"

markdown = """
Yellowport is the second largest city in Sokara. It is mainly a
trading town, and is known for its exotic goods from distant
Ankon-Konu, way to the south.

The Stinking River brings rich deposits of sulphur from the
Lake of the Sea Dragon down to the town, where it is extracted
and stored in the large waterfront warehouses run by the
merchants’ guild. From here, the mineral is exported all over
Harkuna. Unfortunately, all that sulphur has its drawbacks. The
stink is abominable, and much of the city has a yellowish hue.
The river is so full of sulphur that it is virtually useless as a source
of food or of drinking water. However, the demand for sulphur,
especially from the sorcerous guilds, is great.

Politically, much has changed in the past few years. The old
and corrupt king of Sokara, Corin VII, has been deposed and
executed in a military coup. General Grieve Marlock and the
army now control Sokara. The old Council of Yellowport has
been ‘indefinitely dissolved’ and a provost marshal, Marloes
Marlock, the general’s brother, appointed as military governor of
the town.

You can [buy a town house](!purchase-house) in Yellowport for 200 Shards.
Owning a house gives you a place to rest, and to store
equipment.

To leave Yellowport by sea, buy or sell ships and cargo, go
to the harbourmaster.

If you have the codeword Artefact and the Book of the
Seven Sages, you can [turn to 40](!artefact-and-book).

Choose from the following options:

[Seek an audience with the provost marshal](523)

[Visit the market](30)

[Visit the harbourmaster](555)

[Go the merchants’ guild](405)

[Explore the city by day](302)

[Explore the city by night](442)

[Visit your town house](!require-house!300)

[Visit the Gold Dust Tavern](506)

[Visit the temple of Maka](141)

[Visit the temple of Elnir](316)

[Visit the temple of Alvir and Valmir](220)

[Visit the temple of Tyrnai](526)

[Travel north-east towards Venefax](621)

[Head north-west to Trefoille](233)

[Follow the Stinking River north](82)

[Strike out north-west, across country](558)
"""

def on_page(page):
    if codeword_has("assassin"):
        return "50"

    visits = storage_page_get(_track_visits_key)
    if visits == None:
        visits = 0
    visits += 1
    storage_page_set(_track_visits_key, visits)
    if visits == 4:
        return "273"

def on_command(command):
    if command == "artefact-and-book":
        codeword_assert("artefact")
        inventory_assert("bookofthesevensages")
        return "40"

    if command == "purchase-house":
        house_purchase(200)
