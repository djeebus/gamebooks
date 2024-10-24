load("./lib/abilities.star", "get_cunning", abilities_init="init")
load("./lib/bank.star", "cash_get", bank_init="init")

name = "Star smuggler"
start_page = "e001"

def on_start():
    abilities_init()
    bank_init()

def on_page(page):
    page["markdown"] = """
| attributes | value |
| ---------- | ----- |
| cunning | %s |
| cash | %s |

%s
""" % (get_cunning(), cash_get(), page["markdown"])

    if page["page_id"].startswith("r"):
        page["allow_return"] = True
