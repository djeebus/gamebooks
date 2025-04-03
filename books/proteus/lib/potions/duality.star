load("../stats.star", "stat_strength")

_backup_key = "original-strength"
_backup_is_alive_key = "duality-still-alive"


def _on_fight_start():
    storage_set(_backup_key, stat_strength())
    storage_set(_backup_is_alive_key, True)


def _on_damage_taken():
    is_alive = storage_get(_backup_is_alive_key)
    if is_alive == False:
        return 2  # let the default rules deal with it

    strength = stat_strength()
    if strength > 2:
        return 2

    # reset to original strength
    original = storage_get(_backup_key)
    stat_strength(original)
    storage_set(_backup_is_alive_key, False)

    # we already took care of the damage
    return 0

def _on_fight_end():
    if not storage_get(_backup_is_alive_key):
        return

    # reset the user's original strength
    original_strength = storage_get(_backup_key)
    stat_strength(original_strength)


potion_duality = {
    "name": "Duality",
    "description": "An enemy will actually fight someone who only appears to be you. The battle proceeds in the normal way (See Rules for Fighting), but only if the creature wins in that battle will the real you have to fight him.",

    "on_fight_start": _on_fight_start,
    "on_damage_taken": _on_damage_taken,
    "on_fight_end": _on_fight_end,
}
