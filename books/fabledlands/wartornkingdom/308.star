load("../lib/bank.star", "bank_empty")
load("../lib/inventory.star", "inventory_clear")
load("../lib/stamina.star", "stamina_set")

def once():
    bank_empty()
    stamina_set(1)

markdown = """
The ratmen beat you into unconsciousness, and then toss you
down an underground sewer outlet. You are washed up on the
beaches outside Yellowport, where you come to. You have 1
Stamina point left, and the ratmen have taken all the items and
money that you were carrying. Cross them off your Adventure
Sheet. Groggily, you make your way back into the city. [Turn to 10](10).
"""
