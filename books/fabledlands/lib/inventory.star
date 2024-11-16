_inventory_key = '--inventory-key--'

def _get():
    items = storage_get(_inventory_key)
    if items == None:
        items = []
    return items

def _set(items):
    storage_set(_inventory_key, items)

def inventory_add(item):
    storage_push(_inventory_key, item)

def inventory_contains(item_id):
    items = inventory_list()
    for item in items:
        if item['item_id'] == item_id:
            return True
    return False

def inventory_list():
    result = storage_get(_inventory_key)
    if result == None:
        result = []
    return result

def inventory_remove(item_id):
    items = _get()
    for item in items:
        if item["item_id"] == item_id:
            items.remove(item)
            _set(items)
            return

def inventory_assert(item_id):
    items = _get()
    for item in items:
        if item["item_id"] == item_id:
            return
    fail("could not find %s" % item_id)
