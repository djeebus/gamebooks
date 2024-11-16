load("../lib/bank.star", "bank_deposit")

markdown = """
You spin them a tale about how your poor brother, a mercenary
in Grieve Marlock’s personal guard, lost his legs in the fight to
overthrow the old king, and that you have spent all your money
on looking after him. Several of the militia are brought to tears
by your eloquent speech – they end up having a whip-round
among themselves for your brother, and they give you 15
Shards! Chuckling to yourself, you return to the city centre.

[Turn to 10](!commit)
"""


def on_command(command, args):
    if command == 'commit':
        bank_deposit(15)
        return "10"
