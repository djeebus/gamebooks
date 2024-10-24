party_key = "--party-key--"

def get_party():
    party = storage_get(party_key)
    if party is None:
        party = []
    return party

def add_crew(player):
    storage_push(party_key, player)
