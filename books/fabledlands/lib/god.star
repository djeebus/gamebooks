_god_key = "--god--"


def god_get():
    storage_get(_god_key)

def god_set(god):
    storage_set(_god_key, god)

def god_clear():
    storage_remove(_god_key)
