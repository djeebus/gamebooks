load("./lib/inventory.star", "inventory_contains")
load("./lib/stats.star", "stats_init", "stat_courage", "stat_dexterity", "stat_strength")

start_page = "start"

def on_start():
    stats_init()

def on_page(page):
    page["markdown"] = """
%s

Dexterity: %d
Strength: %d
Courage: %d    
""" % (page["markdown"], stat_dexterity(), stat_strength(), stat_courage())

    page.setdefault("clear_history", True)


def _must_have_item(item_id, page_id):
    if not inventory_contains(item_id):
        fail("inventory does not contain %s" % item_id)

    return page_id


def on_command(cmd, args):
    if cmd == "must-have-item":
        return _must_have_item(*args)
