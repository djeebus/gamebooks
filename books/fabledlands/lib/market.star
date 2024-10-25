load("../lib/bank.star", bank_balance="balance", bank_deposit="deposit", bank_assert_min_balance="assert_min_balance")
load("../lib/inventory.star", inventory_add="add", inventory_remove="remove")

def render(name, items, *item_ids):
    rows = [_render_item(items, item_id) for item_id in item_ids]

    return """
| %s | To buy | To sell |
| ------- | ----- | ------- |
%s  
    """ % (name, '\n'.join(rows))


def _create_link(item_id, action, cost):
    if not cost:
        return "-"

    return "[%d Shards](!%s!%s)" % (cost, action, item_id)


def _render_item(items, item_id):
    item = items[item_id]

    buy = _create_link(item_id, "buy", item.get("buy"))
    sell = _create_link(item_id, "sell", item.get("sell"))

    return "| %s | %s | %s |" % (item["item"]["label"], buy, sell)

def on_command(command, items):
    if command.startswith("buy!"):
        item_id = command.lstrip("buy!")
        info = items[item_id]
        cost = info["buy"]
        item = info["item"]
        bank_assert_min_balance(cost)
        inventory_add(item)
        return True

    if command.startswith("sell!"):
        item_id = command.lstrip("sell!")
        info = items[item_id]
        value = info["sell"]
        item = info["item"]
        inventory_remove(item)
        bank_deposit(value)
        return True

    return False