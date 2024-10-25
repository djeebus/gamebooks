_codewords_key = '--codewords--'


def codeword_list():
    words = storage_get(_codewords_key)
    if words == None:
        words = list()
    return words

def codeword_add(codeword):
    words = codeword_list()
    if codeword not in words:
        words.append(codeword)
    storage_set(_codewords_key, words)


def codeword_remove(codeword):
    words = codeword_list()
    words = [w for w in codeword if w != codeword]
    storage_set(_codewords_key, words)


def codeword_assert(codeword):
    words = codeword_list()
    if codeword not in words:
        fail("player does not have codeword")
