cunning_key = "--cunning--"

def init():
    cunning_score = dice_roll(1, 6)
    storage_set(cunning_key, cunning_score)

def get_cunning():
    return storage_get(cunning_key)
