_key = "strength-damage-counter"


def _on_fight_start():
    storage_page_set(_key, 0)


def _on_damage_given(_):
    index = storage_page_get(_key)
    if index < 2:
        index += 1
        storage_page_set(_key, index)
        return 4

    return 2


potion_strength = {
    "name": "Power",
    "description": "Each blow in any one battle will cause four Strength points to be lost by the enemy for the first two rounds of fighting (See Rules for Fighting).",

    "on_fight_start": _on_fight_start,
    "on_damage_given": _on_damage_given,
}
