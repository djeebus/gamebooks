load("./stats.star", "stat_strength")

_meals_key = "meals"


def meals_add():
    count = meals_get()
    meals_set(count + 1)


def meals_get():
    count = storage_get(_meals_key)
    if count == None:
        count = 0
    return count


def meals_set(count):
    storage_set(_meals_key, count)


def meals_eat():
    count = meals_get()
    if count <= 0:
        fail("no meals")

    strength = stat_strength()
    storage_set(_meals_key, count - 1)
    stat_strength(_meals_key, strength + 5)
