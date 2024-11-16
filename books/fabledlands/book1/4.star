load("../lib/bank.star", "bank_deposit")
load("../lib/inventory.star", "inventory_add")

def once():
    bank_deposit(100)
    inventory_add({
        "id": "trident",
        "label": "Trident (COMBAT +1)",
        "attributes": {
            "combat-mode": 1,
        },
    })

markdown = """
The priests of Alvir and Valmir are overjoyed that you have
returned the golden net. The high priest rewards you with 100
Shards and a magic weapon, a rune-engraved trident.

[Next page](220)
"""
