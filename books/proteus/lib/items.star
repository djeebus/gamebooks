load("./stats.star", "stat_strength")
load("./potions/duality.star", "potion_duality")
load("./potions/elusiveness.star", "potion_elusiveness")
load("./potions/fear.star", "potion_fear")
load("./potions/invincibility.star", "potion_invincibility")
load("./potions/strength.star", "potion_strength")

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
    "invincibility-potion": potion_invincibility,
    "flying-potion": {
        "name": "Flying",
        "description": "You will float over any natural obstacle.",
    },
    "power-potion": potion_strength,
    "calm-potion": {
        "name": "Calm",
        "description": "Will restore your Courage points to their Initial level. ",
    },
    "fear-potion": potion_fear,
    "intuition-potion": {
        "name": "Intuition",
        "description": "You will know the answer to a question, without necessarily knowing why.",
    },
    "duality-potion": potion_duality,
    "elusiveness-potion": potion_elusiveness,
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
