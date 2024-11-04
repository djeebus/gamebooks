_cash_key = "--cash--"


def bank_assert_min_balance(amount):
    cash = bank_balance()
    if cash < amount:
        fail("not enough money in bank (%d < %d)" % (cash, amount))

def bank_balance():
    cash = storage_get(_cash_key)
    if cash == None:
        cash = 0
    return cash

def bank_deposit(amount):
    cash = bank_balance()
    cash += amount
    storage_set(_cash_key, cash)

def bank_withdraw(amount):
    bank_assert_min_balance(amount)
    cash = bank_balance()
    cash -= amount
    storage_set(_cash_key, cash)
