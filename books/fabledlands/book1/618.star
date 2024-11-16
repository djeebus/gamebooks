load("../lib/bank.star", "bank_withdraw")
load("../lib/god.star", "god_set")

markdown = """
Becoming an initiate of Lacuna gives you the benefit of paying
less for blessings and other services the temple can offer. It costs
30 Shards to become an initiate. You cannot do this if you are
already an initiate of another temple. 

[Choose to become an initiate](!purchase) - Costs 30 shards
If you choose to become an initiate, write Lacuna in the God box on your Adventure
Sheet – and remember to cross off the 30 Shards.

[Return](544)
"""

def on_command(command):
    if command != "purchase":
        return

    bank_withdraw(30)
    god_set("Lacuna")
