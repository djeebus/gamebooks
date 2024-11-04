load("../lib/bank.star", "bank_deposit")
load("../lib/inventory.star", "inventory_add")
load("../lib/items.star",
     "weapon3", "armor5", "magical2", "mandolin2", "goldcompass",
     "magiclockpicks", "silverholysymbol")

items = [
    {
        "item": weapon3,
    },
    {
        "item": armor5,
    },
    {
        "item": magical2,
    },
    {
        "cash": 500,
    },
    {
        "item": mandolin2,
    },
    {
        "item": goldcompass,
    },
    {
        "item": magiclockpicks,
    },
    {
        "item": silverholysymbol,
    },
]


def _label(row):
    if "item" in row:
        return row["item"]["label"]
    if "cash" in row:
        return "%d Shards" % (row["cash"])
    fail("unknown item")


def _key(index):
    return "key:%d" % index


def _select(index):
    key = _key(index)
    storage_page_set(key, True)


def _selected(index):
    key = _key(index)
    return storage_page_get(key) == True


def _unselect(index):
    key = _key(index)
    storage_page_remove(key)


rows = [
    "- [%s] [%s](!%s!%d)" % (
        "x" if _selected(index) else " ",
        _label(row),
        "unselect" if _selected(index) else "select",
        index,
    ) for index, row in enumerate(items)
]

markdown = """
If there is a tick in the box, turn to 251 immediately. If not, put
a tick there now and read on.

If you have the codeword Avenge, turn to 648 immediately.
Otherwise read on.

You remain quiet as a mouse, behind a pile of coins. After a
long wait, the sea dragon slithers into the water, and swims out
on some errand. You have some time to loot the hoard. You
may choose up to three of the following treasures:

%s

After you have taken the third treasure, you hear the sea dragon
returning. Quickly you climb up through the hole in the roof
on to an island in the middle of the lake. From there you
manage to get a lift on a passing boat, and make it safely to
Cadmium village. 

[Next page](!commit)
""" % ("\n".join(rows))


def _sum(items):
    total = 0
    for item in items:
        total += item
    return total


def on_command(command):
    parts = command.split('!')
    cmd = parts[0]
    args = parts[1:]

    if cmd == "select":
        index = int(args[0])
        _select(index)
        return

    if cmd == "unselect":
        index = int(args[0])
        _unselect(index)
        return

    if cmd == "commit":
        selected = [
            1 if _selected(index) else 0
            for index, item in enumerate(items)
        ]
        total = _sum(selected)
        if total != 3:
            fail("%d != 3" % total)

        for index, item in enumerate(items):
            if not _selected(index):
                continue

            if "item" in item:
                inventory_add(item["item"])
            elif "cash" in item:
                bank_deposit(item["cash"])

        return "135"

    fail("unknown command: %s" % command)
