load("./lib/abilities.star", "cunning_get", abilities_init="init")
load("./lib/bank.star", "cash_get", bank_init="init")
load("./lib/location.star", "location_get", location_init="init")

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
| location | %s |

%s

[source material](https://dwarfstar.brainiac.com/ds_starsmuggler.html)
""" % (cunning_get(), cash_get(), location_get(), page["markdown"])

    if page["page_id"].startswith("r"):
        page["allow_return"] = True
