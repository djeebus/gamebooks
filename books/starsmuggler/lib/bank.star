cash_key = '--cash-key--'

def init():
    initial_cash = dice_roll(1, 6) * 100 + 150
    storage_set(cash_key, initial_cash)

def get_cash():
    return storage_get(cash_key)
