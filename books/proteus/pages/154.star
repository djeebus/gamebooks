load("../lib/fight.star", "fight_render", "fight_build_command")

markdown = """
He laughs. he is completely unaffected.
“Your magic is too weak against such a powerful 
Wizard as I am,” he bellows. You parry a
sudden blow from the mighty axe with your
shield, but are knocked backwards. You realise
that he is enormously strong. Now it is a fight
to the death:

%s
""" % fight_render("BELENGHAST", 14, 22)

on_command = fight_build_command({
    "goto": "198",
})
