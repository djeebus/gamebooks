markdown = """
The ZOMBIE backs away, eyeing your silver
lance warily. You close for the attack. Even
though you have the lance, the ZOMBIE is a
formidable enemy.

%s
""" % (fight_render("ZOMBIE", 11, 20))

_blows_key = "blows"

def _on_damage_dealt(info):
    blows = storage_page_get(_blows_key)
    blows = 0 if blows == None else blows
    blows += 1
    if blows == 3:
        return info["strength"]["current"]
    storage_page_set(_blows_key, blows)


on_command = fight_build_command({
    "on_damage_dealt": _on_damage_dealt,
    "goto": 116,
})
