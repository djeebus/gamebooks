keys = ["charisma", "combat", "magic", "sanctity", "scouting", "thievery"]

def _key(ability):
    return "--ability-%s--" % ability

def ability_get(ability):
    key = _key(ability)
    return storage_get(key)

def set(ability, score):
    key = _key(ability)
    storage_set(key, score)
