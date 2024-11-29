_key = "invincibility-damage-counter"


def _on_fight_start():
    storage_page_set(_key, 0)


def _on_damage_taken(_):
    index = storage_page_get(_key)
    if index < 3:
        index += 1
        storage_page_set(_key, index)
        return 1

    return 2


potion_invincibility = {
    "name": "Invincibility",
    "description": "Not quite as powerful as it sounds! It will totally protect you against any creature only for the first three rounds of fighting (See Rules for Fighting). This means you will lose no Strength points for the first three times you are wounded in that battle.",

    "on_fight_start": _on_fight_start,
    "on_damage_taken": _on_damage_taken,
}
