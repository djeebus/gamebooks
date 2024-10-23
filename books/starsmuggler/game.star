load("./lib/abilities.star", "get_cunning", abilities_init="init")
load("./lib/bank.star", "get_cash", bank_init="init")
name = "Star smuggler"
start_page = "e001"

def on_start():
    abilities_init()

def on_page(page):
    page["markdown"] = """
| attributes | value |
| ---------- | ----- |
| cunning | %s |
| cash | %s |

%s
""" % (get_cunning(), get_cash(), page["markdown"])
