load("./stats.star", "stat_strength")

def _elusiveness_potion_on_fight_start():
    storage_page_set("damage-counter", 0)

def _elusiveness_potion_on_damage_taken(_):
    index = storage_page_get("damage-counter")
    if index < 3:
        index += 1
        storage_page_set("damage-counter", index)
        return 1

    return 2

_double_strength_key = "--double-strength--"

def _duality_potion_on_damage_taken(enemy):
    double_strength = storage_get(_double_strength_key)
    if double_strength == None:
        double_strength = stat_strength()
        storage_set(_double_strength_key, double_strength)

    damage = 2
    if enemy["strength"] <= 2:

    # if damage would drop to 0
    # if double hasn't died yet
    # reset health to original
    # mark double as dead

def _duality_potion_on_fight_end():
    # if double is not dead
    # reset health to original

all_items = {
    "truth-seeking-potion": {
        "name": "Truth-seeking potion",
        "description": "You will know when someone is lying to you.",
    },
    "searching-potion": {
        "name": "Searching",
        "description": "Will take you in the right direction when you are faced with a choice.",
    },
    "transparency-potion": {
        "name": "Transparency",
        "description": "Enables you to disappear, like a ghost.",
    },
    "invincibility-potion": {
        "name": "Invincibility",
        "description": "Not quite as powerful as it sounds! It will totally protect you against any creature only for the first three rounds of fighting (See Rules for Fighting). This means you will lose no Strength points for the first three times you are wounded in that battle.",
    },
    "flying-potion": {
        "name": "Flying",
        "description": "You will float over any natural obstacle.",
    },
    "power-potion": {
        "name": "Power",
        "description": "Each blow in any one battle will cause four Strength points to be lost by the enemy for the first two rounds of fighting (See Rules for Fighting).",
    },
    "calm-potion": {
        "name": "Calm",
        "description": "Will restore your Courage points to their Initial level. ",
    },
    "fear-potion": {
        "name": "Fear",
        "description": "Any enemy will lose three Dexterity points for that particular battle (See Rules for Fighting).",
    },
    "intuition-potion": {
        "name": "Intuition",
        "description": "You will know the answer to a question, without necessarily knowing why.",
    },
    "duality-potion": {
        "name": "Duality",
        "description": "An enemy will actually fight someone who only appears to be you. The battle proceeds in the normal way (See Rules for Fighting), but only if the creature wins in that battle will the real you have to fight him.",
        "on_damage_taken": _duality_potion_on_damage_taken,
        "on_fight_end": _duality_potion_on_fight_end,
    },
    "elusiveness-potion": {
        "name": "Elusiveness",
        "description": "This makes you difficult for an enemy to hit. You will lose only 1 Strength point for the first three rounds of that battle (See Rules for Fighting).",
        "on_fight_start": _elusiveness_potion_on_fight_start,
        "on_damage_taken": _elusiveness_potion_on_damage_taken,
    },
    "madness-potion": {
        "name": "Madness",
        "description": "An enemy’s actions become completely wild and unpredictable.",
    },
    "revitalization-potion": {
        "name": "Revitalization",
        "description": "Restores Strength points to their Initial level (See Rules for Fighting).",
    },

    "amulet-of-stone": {
        "name": "Amulet of Stone",
    },
}
