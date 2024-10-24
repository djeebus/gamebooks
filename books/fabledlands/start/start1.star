load("../lib/character.star", "character_render", "character_select")

title = "Starting characters"

_common = {
    "rank": 1,
    "stamina": 9,
    "money": 16,
}

_characters = {
    "liana": dict(
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
        **_common,
    ),
    "andriel": dict(
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
        **_common,
    ),
    "chalor": dict(
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
        **_common,
    ),
    "marana": dict(
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
        **_common,
    ),
}

markdown = """
You can create your own character, or pick one from the following.
Transfer the details of the character you have chosen to the Adventure Sheet.

## [Liana The Swift](!liana)

%s
 
## [Andriel The Hammer](!andriel)

%s

## [Chalor The Exiled One](!chalor)

%s

## [Marana Fireheart](!marana)

%s
 
""" % (
    character_render(_characters["liana"]),
    character_render(_characters["andriel"]),
    character_render(_characters["chalor"]),
    character_render(_characters["marana"]),

)

can_submit = True

def on_command(name):
    character_select(_characters[name])
    return "start2"
