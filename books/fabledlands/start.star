load("./lib/abilities.star", "ability_set", ability_keys="keys")
load("./lib/bank.star", "bank_deposit")
load("./lib/character.star", "defense_set", "name_set", "profession_set", "rank_set")
load("./lib/stamina.star", "stamina_set")


def character_render(c_id, c):
    return """
    
## [%s](!%s)

|   |    |
|---|----|
| Profession | %s |
| Defence | %s |
|     |     |
| Charisma | %d |
| Combat | %d |
| Magic | %d |
| Sanctity | %d |
| Scouting | %d |
| Thievery | %d |
    """ % (
        c["name"],
        c_id,
        c["profession"],
        c["defence"],
        c["charisma"],
        c["combat"],
        c["magic"],
        c["sanctity"],
        c["scouting"],
        c["thievery"],
    )


_characters = {
    "liana": dict(
        name="Liana The Swift",
        profession="Wayfarer",
        defence=7,
        charisma=2,
        combat=5,
        magic=2,
        sanctity=3,
        scouting=6,
        thievery=4,
        posessions=[
            "spear",
            "leather jerkin (Defence +1)",
            "map",
        ],
    ),
    "andriel": dict(
        name="Andriel The Hammer",
        profession="Warrior",
        defence=8,
        charisma=3,
        combat=6,
        magic=2,
        sanctity=4,
        scouting=3,
        thievery=2,
        posessions=[
            "battle-axe",
            "leather jerkin (Defence +1)",
            "map",
        ],
    ),
    "chalor": dict(
        name="Chalor The Exiled One",
        profession="Mage",
        defence=4,
        charisma=2,
        combat=2,
        magic=6,
        sanctity=1,
        scouting=5,
        thievery=3,
        posessions=[
            "staff",
            "leather jerkin (Defence +1)",
            "map",
        ],
    ),
    "marana": dict(
        name="Marana Fireheart",
        profession="Rogue",
        defence=6,
        charisma=5,
        combat=4,
        magic=4,
        sanctity=1,
        scouting=2,
        thievery=6,
        posessions=[
            "sword",
            "leather jerkin (Defence +1)",
            "map",
        ],
    ),
}

markdown = """
# Starting characters

You can create your own character, or pick one from the following.
Transfer the details of the character you have chosen to the Adventure Sheet.

%s

""" % "\n\n".join([character_render(key, value) for (key, value) in _characters.items()])

can_submit = True


def on_command(name):
    character = _characters[name]
    for key in ability_keys:
        ability_set(key, character[key])

    name_set(character["name"])
    profession_set(character["profession"])
    rank_set(1)
    defense_set(character["defence"])
    bank_deposit(10)
    stamina_set(10)

    return "book1/1"
