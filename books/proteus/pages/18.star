load("../lib/fight.star", "fight_render", "fight_build_command")

_damage_counter = "damage-counter"


def once():
    storage_set_page(_damage_counter, 0)


def _damage_calculator():
    index = storage_get_page(_damage_counter)
    if index < 3:
        index += 1
        storage_set_page(_damage_counter, index)
        return 1

    return 2


markdown = """
The potion works – you are now very
difficult to wound. You will lose only 1
Strength point, instead of 2, for the first three
successful attacks against you. After that, the
potion wears off, and each successful attack
against you by the SCRAFE will cause the
usual 2 points of damage.

%s

If you win, turn to 152.
""" % (fight_render("SCRAFE", 10, 12))


on_command = fight_build_command({
    "calculate_damage_taken": _damage_calculator,
    "goto": "152",
})
