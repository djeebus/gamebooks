load("../lib/codewords.star", "codeword_has")
load("../lib/inventory.star", "inventory_contains")

markdown = """
The soldier recognizes you. He bows and says, ‘Welcome,
my lord. I will take you see King Nergan.’

He leads you to Nergan’s mountain stockade, where the
king greets you warmly.

‘Ah, my local champion! It is always a pleasure to see you.
However, I was hoping you had spoken with General Beladai of
the allied army – we need that citadel. Now go. That is a royal
command!’

You leave, climbing down to the foothills of the mountains.
[Next page](474).
"""

def on_page(page):
    if inventory_contains("coded missive"):
        return "676"

    if codeword_has("deliver"):
        return "98"
