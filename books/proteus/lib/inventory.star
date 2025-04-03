
_inventory_storage_key = "inventory_item_ids"


def inventory_add(item_id):
    item_ids = storage_get(_inventory_storage_key) or []
    item_ids = list(item_ids)  # make sure it's a list
    item_ids.extend([item_id])
    storage_set(_inventory_storage_key, item_ids)


def inventory_contains(item_id):
    item_ids = storage_get(_inventory_storage_key) or []
    return item_id in item_ids
