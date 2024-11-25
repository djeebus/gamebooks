_dexterity_key = "dexterity"
_initial_dexterity_key = "initial-dexterity"

_strength_key = "strength"
_initial_strength_key = "initial-strength"

_courage_key = "courage"
_initial_courage_key = "initial-courage"


def stats_init():
    dexterity = dice_roll(1, 6) + 8
    storage_set(_dexterity_key, dexterity)
    storage_set(_initial_dexterity_key, dexterity)
    log("calculated initial dexterity: %d" % dexterity)

    strength = dice_roll(2, 6) + 15
    storage_set(_strength_key, strength)
    storage_set(_initial_strength_key, strength)
    log("calculated initial strength: %d" % strength)

    courage = dice_roll(1, 6) + 6
    storage_set(_courage_key, courage)
    storage_set(_initial_courage_key, courage)
    log("calculated initial courage: %d" % courage)


def _get_or_set(value_key, initial_key, new_value):
    if new_value == None:
        return storage_get(value_key)

    init_value = storage_get(initial_key)
    if new_value > init_value:
        new_value = init_value
    storage_set(value_key, new_value)
    log("%s = %d" % (value_key, new_value))
    return new_value


def _add(current_key, initial_key, new_value):
    current = storage_get(current_key)
    current += new_value
    if current <= 0:
        fail("%s has dropped to zero; game over!" % current_key)

    initial = storage_get(initial_key)
    if current > initial:
        current = initial
    storage_set(current_key, current)

    if new_value < 0:
        log("you have lost %d %s" % (new_value, current_key))
    else:
        log("you have gained %d %s" % (new_value, current_key))


def stat_dexterity(value = None):
    return _get_or_set(_dexterity_key, _initial_dexterity_key, value)


def stat_strength(value = None):
    return _get_or_set(_strength_key, _initial_strength_key, value)


def strength_add(value):
    return _add(_strength_key, _initial_strength_key, value)


def stat_courage(value = None):
    return _get_or_set(_courage_key, _initial_courage_key, value)


def courage_add(value):
    return _add(_courage_key, _initial_courage_key, value)
