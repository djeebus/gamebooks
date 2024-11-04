_codewords_key = '--codewords--'

_all_words = [
    "acid", "anvil", "afraid", "apache", "ague", "appease", "aid", "apple", "aklar", "ark", "alissia", "armour",
    "almanac", "artefact", "aloft", "artery", "altitude", "ashen", "altruist", "aspen", "ambuscade", "assassin",
    "amcha", "assault", "amends", "assist", "anchor", "attar", "anger", "avenge", "animal", "axe", "anthem", "azure",
]


def _clean(codeword):
    return codeword.strip().lower()


def _get():
    words = storage_get(_codewords_key)
    if words == None:
        words = dict()
    return words


def _set(words):
    storage_set(_codewords_key, words)


def codeword_all():
    words = _get()
    for word in _all_words:
        words.setdefault(word, False)
    return words


def codeword_has(codeword):
    words = _get()
    return codeword in words

def codeword_add(codeword):
    codeword = _clean(codeword)
    words = _get()
    words[codeword] = True
    _set(words)


def codeword_remove(codeword):
    codeword = _clean(codeword)
    words = _get()
    words.pop(codeword)
    _set(words)


def codeword_assert(codeword):
    codeword = _clean(codeword)
    words = _get()
    if codeword not in words:
        fail("player does not have codeword")
