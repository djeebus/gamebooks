load("../lib/bank.star", "bank_withdraw")
load("../lib/god.star", "god_clear")

markdown = """
To renounce the worship of Lacuna, you must pay 40 Shards in
compensation to the shrine.

‘If you forsake the love of the goddess, you will never
survive the rigours of the wilderness,’ warns the priestess.

[Renounce the worship of Lacuna](!forsake) - 40 Shards.

[Return](195)
"""

def on_command(command):
    if command != "forsake":
        fail("unknown command")

    bank_withdraw(40)
    god_clear()
