character_key = "--character--"

def character_render(c):
    return """
|   |    |
|---|----|
| Rank | %s |
| Profession | %s |
| Stamina | %s |
| Defence | %s |
| Money | %d shards |
|     |     |
| Charisma | %d |
| Combat | %d |
| Magic | %d |
| Sanctity | %d |
| Scouting | %d |
| Thievery | %d |
    """ % (
        c["rank"],
        c["profession"],
        c["stamina"],
        c["defence"],
        c["money"],
        c["charisma"],
        c["combat"],
        c["magic"],
        c["sanctity"],
        c["scouting"],
        c["thievery"],
)

def character_select(c):
    storage_set(character_key, c)
