health_key = "--health--"
default_health = 10


def get_max():
    return 10


def get():
    health = storage_get(health_key)
    if health == None:
        health = default_health
    return health


def set(value):
    return storage_set(health_key, value)


def add(value):
    health = get()
    max_stamina = get_max()
    if health < max_stamina:
        health += value
        set(health)
