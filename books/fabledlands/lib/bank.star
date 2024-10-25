_cash_key = "--cash--"


def balance():
    cash = storage_get(_cash_key)
    if cash == None:
        cash = 0
    return cash

def withdraw(amount):
    assert_min_balance(amount)
    cash = balance()
    cash -= amount
    storage_set(_cash_key, cash)

def deposit(amount):
    cash = balance()
    cash += amount
    storage_set(_cash_key, cash)

def assert_min_balance(amount):
    cash = balance()
    if cash < amount:
        fail("not enough money in bank (%d < %d)" % (cash, amount))
