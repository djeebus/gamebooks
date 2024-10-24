cunning_key = "--cunning--"

def init():
    cunning_score = dice_roll(1, 6)
    storage_set(cunning_key, cunning_score)

def get_cunning():
    return storage_get(cunning_key)

def make_cunning_roll():
    cunning_score = get_cunning()
    if cunning_score == None:
        fail("cunning score is missing!")
    roll = dice_roll(1, 6)
    success = roll <= cunning_score
    return success
