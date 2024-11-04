health_key = "--health--"
default_health = 10


def stamina_get_max():
    return 10


def stamina_get():
    health = storage_get(health_key)
    if health == None:
        health = default_health
    return health


def stamina_set(value):
    return storage_set(health_key, value)


def stamina_add(value):
    health = stamina_get()
    health += value

    if health <= 0:
        fail("you are dead. sad times.")

    max_stamina = stamina_get_max()
    if health < max_stamina:
        health = max_stamina
        stamina_set(health)
