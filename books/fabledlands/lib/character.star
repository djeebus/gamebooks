name_key = "--name--"

def name_get():
    return storage_get(name_key)

def name_set(value):
    storage_set(name_key, value)

profession_key = "--profession--"

def profession_get():
    return storage_get(profession_key)

def profession_set(value):
    storage_set(profession_key, value)

rank_key = '--rank--'

def rank_get():
    return storage_get(rank_key)

def rank_set(value):
    storage_set(rank_key, value)

defense_key = '--defense--'

def defense_get():
    return storage_get(defense_key)

def defense_set(value):
    storage_set(defense_key, value)
