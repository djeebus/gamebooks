markdown = """
This is another test page.

[markdown](markdown)
"""

def on_page(page):
    return {
        "markdown": "page.on_page\n\n" + page["markdown"]
    }
