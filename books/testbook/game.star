_init_key = "init"


def on_start():
    storage_set(_init_key, True)


def on_page(page):
    has_initted = storage_get(_init_key)

    page["markdown"] += "\n\nbook.on_page: global storage works: %s" % ("yes" if has_initted else "no")

    return {
        "markdown": page["markdown"] + "\n\nbook.on_page: global storage works: %s" % ("yes" if has_initted else "no"),
    }


start_page = "tests"
