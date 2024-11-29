_key = "elusiveness-damage-counter"


def _on_fight_start():
    storage_page_set(_key, 0)


def _on_damage_taken(_):
    index = storage_page_get(_key)
    if index < 3:
        index += 1
        storage_page_set(_key, index)
        return 1

    return 2


potion_elusiveness = {
    "name": "Elusiveness",
    "description": "This makes you difficult for an enemy to hit. You will lose only 1 Strength point for the first three rounds of that battle (See Rules for Fighting).",

    "on_fight_start": _on_fight_start,
    "on_damage_taken": _on_damage_taken,
}
