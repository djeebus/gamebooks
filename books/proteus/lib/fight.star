load("./stats.star", "strength_add", "stat_dexterity")

_fight_command = "fight"
_strength_key_fmt = "enemy:%d"


def _fight_round(enemies, options):
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
        on_damage_dealt = options.get("on_damage_dealt")
        damage_dealt = on_damage_dealt(enemy) if on_damage_dealt != None else None
        damage_dealt = 2 if damage_dealt == None else damage_dealt
        log("you have hit the %s for %d damage" % (enemy["name"], damage_dealt))
        enemy_strength -= damage_dealt
        storage_page_set(_strength_key_fmt % index, enemy_strength)

        return enemy_strength <= 0 if is_last_enemy else False

    if self_power < enemy_power:
        calculator = options.get("calculate_damage_taken")
        damage_taken = calculator if calculator != None else 2
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
