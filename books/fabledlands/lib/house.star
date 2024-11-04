load("../lib/bank.star", "bank_withdraw")

_key = "house-purchase"

def house_owned():
    value = storage_page_get(_key)
    return value == True


def house_purchase(price):
    if house_owned():
        fail("house already owned")

    bank_withdraw(price)
    storage_page_set(_key, True)


def house_assert():
    if house_owned() == False:
        fail("you do not own a house here")
