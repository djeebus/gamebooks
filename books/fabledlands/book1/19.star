load("../lib/inventory.star", "inventory_add")
load("../lib/items.star", "armor6")

_tick_key = 'ticks'

def _get_ticks():
    num = storage_page_get(_tick_key)
    if num == None:
        num = 0
    return num

def _inc_tick(max_tick):
    num = _get_ticks()
    if num < max_tick:
        num += 1
        storage_page_set(_tick_key, num)
    return num


def _maybe(t_or_f, label):
    if t_or_f:
        return '[' + label + '](!gain-codeword!anvil)'

    return '~label~'


markdown = """
Put a tick in an empty box. If all three boxes are now ticked,
%s.

The Dragon Knights are impressed with your combat skills.
Your opponent comes round, ruefully rubbing his neck.
Grudgingly, he admits to your superior skill and hands you his
weapon and armour. You get an ordinary sword and a suit of
heavy plate (Defence +6). You take your leave, [turn to 276](!next).
""" % (_maybe(_get_ticks() == 3, "gain the codeword Anvil"))

def on_command(command, *args):
    if command == "next":
        inventory_add(armor6)
        return "276"
