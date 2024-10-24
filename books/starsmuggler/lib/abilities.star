cunning_key = "--cunning--"

def init():
    cunning_score = dice_roll(1, 6)
    storage_set(cunning_key, cunning_score)

def cunning_get():
    return storage_get(cunning_key)

def cunning_roll():
    cunning_score = cunning_get()
    if cunning_score == None:
        fail("cunning score is missing!")
    roll = dice_roll(1, 6)
    success = roll <= cunning_score
    return success
