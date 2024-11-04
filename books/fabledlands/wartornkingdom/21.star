load("../lib/bank.star", "bank_deposit", "bank_balance", "bank_withdraw")
load("../lib/stamina.star", "stamina_set")

markdown = """
While making your way through the back streets of the poor
quarter you are set upon a knife-wielding thug, who is intent of
relieving you of your purse.

If you don’t want to fight him, you can try a [CHARISMA roll
at a Difficulty of 8](!roll!charisma!8!10!) to try to talk your way out of this unpleasant
situation. If you succeed, the thug leaves, confused by your
rhetoric (turn to 10 and choose again.) Otherwise, you must
fight him.

[Thug, COMBAT 4, Defence 7, Stamina 6](!fight!4!7!6!success!fail)

If you defeat him, you find 15 Shards on his body. If you are
defeated, you are stunned into unconsciousness. You come
round with 1 Stamina point, and he has robbed you of 50 Shards
(or of all your money if you have less than 50 Shards). Turn to 10.
"""

def on_command(command):
    if command == "success":
        bank_deposit(15)
        return "10"

    if command == "fail":
        stamina_set(1)
        balance = bank_balance()
        loss = balance if balance < 50 else 50
        bank_withdraw(loss)
        return "10"
