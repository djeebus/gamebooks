load("../lib/bank.star", bank_withdraw="withdraw")
load("../lib/stamina.star", stamina_add="add")

markdown = """
The Green Man Inn costs you [1 Shard a day](!rest). Each day you
spend here, you can recover 1 Stamina point if injured, up to
the limit of your normal unwounded Stamina score.

[Return](195)
"""

def on_command(command):
    if command == "rest":
        bank_withdraw(1)
        stamina_add(1)
        return

    fail("unknown command: %s" % command)
