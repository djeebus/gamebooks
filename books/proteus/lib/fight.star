load("./items.star", "all_items")
load("./stats.star", "strength_add", "stat_dexterity")

_fight_command = "fight"
_fight_start_key = "fight-init"
_strength_key_fmt = "enemy:%d"

def _default_on_fight_start():
    pass

def _default_on_fight_end():
    pass

def _default_on_damage_given(_):
    return 2

def _default_on_damage_taken(_):
    return 2


def _fight_round(enemies, options):
    callbacks = {
        "on_fight_start": _default_on_fight_start,
        "on_damage_given": _default_on_damage_given,
        "on_damage_taken": _default_on_damage_taken,
        "on_fight_end": _default_on_fight_end,
    }

    item_id = options.get("item_id")
    if item_id:
        item = all_items[item_id]
        for key, val in item.items():
            if key in callbacks:
                callbacks[key] = val

    init = storage_page_get(_fight_start_key)
    if init != True:
        callbacks["on_fight_start"]()
        storage_page_set(_fight_start_key, True)

    enemy = None
    for index, enemy in enumerate(enemies):
        enemy_strength = storage_page_get(_strength_key_fmt % index)
        enemy_strength = enemy_strength if enemy_strength != None else enemy["strength"]
        if enemy_strength > 0:
            break

    if enemy == None:
        return fail("no more enemies to fight!")

    is_last_enemy = index == len(enemies) - 1

    enemy_dexterity = enemy["dexterity"]

    self_power = dice_roll(2, 6) + stat_dexterity()
    enemy_power = dice_roll(2, 6) + enemy_dexterity

    log("your power = %d" % self_power)
    log("enemy power = %d" % enemy_power)

    if self_power > enemy_power:
        damage_dealt = callbacks["on_damage_given"](enemy)
        log("you have hit the %s for %d damage" % (enemy["name"], damage_dealt))
        enemy_strength -= damage_dealt
        storage_page_set(_strength_key_fmt % index, enemy_strength)

        if is_last_enemy:
            if enemy_strength <= 0:
                callbacks["on_fight_end"]()
                return True

        return False

    if self_power < enemy_power:
        damage_taken = callbacks["on_damage_taken"](enemy)
        strength_add(-damage_taken)

_enemy_key = "enemy-data"

def fight_render(name, dexterity, strength):
    return fight_multiple_render([{
        "name": name,
        "dexterity": dexterity,
        "strength": strength,
    }])


def _render_row(index, row):
    current = storage_page_get(_strength_key_fmt % index)
    if current == None:
        current = row["strength"]

    return "| %s   |        %d | %s / %s  |" % (
        row["name"],
        row["dexterity"],
        current,
        row["strength"],
    )


def fight_multiple_render(enemies):
    fight_data = storage_page_get(_enemy_key)
    if fight_data == None:
        log("setting initial data")
        fight_data = enemies
        storage_page_set(_enemy_key, fight_data)

    rows = [
        _render_row(index, row)
        for index, row in enumerate(enemies)
    ]

    return """
| name | dexterity | strength |
| ---- | --------- | -------- |
%s

[Fight](!%s)    
""" % (
        "\n".join(rows),
        _fight_command,
    )


def fight_build_command(options):
    def on_command(cmd, args):
        if cmd == _fight_command:
            enemies = storage_page_get(_enemy_key)
            result = _fight_round(enemies, options)
            if result == True:
                return options["goto"]
    return on_command
