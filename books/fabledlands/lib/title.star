_titles = [
    'Protector of Sokara',
]

_title_key = 'titles'


def _clean(key):
    return key.lower().strip()


def _get():
    titles = storage_get(_title_key)
    if titles == None:
        titles = []
    return titles


def title_has(title):
    title = _clean(title)
    titles = _get()
    return title in titles


def title_assert(title):
    if title_has(title) == False:
        fail("title not posessed")
