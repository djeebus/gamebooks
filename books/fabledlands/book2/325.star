load("../lib/character.star", "rank_get")

_roll1_key = 'die-roll-1'
_roll2_key = 'die-roll-2'

def on_page(page):
    roll_1 = storage_page_get(_roll1_key)
    if roll_1 == None:
        roll_1 = dice_roll(2, 6)
        storage_page_set(_roll1_key, roll_1)

    roll_2 = storage_page_get(_roll2_key)
    if roll_2 == None:
        roll_2 = dice_roll(1, 6)
        storage_page_set(_roll2_key, roll_2)

    rank = rank_get()
    if rank < roll_1:
        page["markdown"] = """
    Your ship, crew and cargo are lost to the deep, dark sea. Cross
    them off your Adventure Sheet. Your only thought now is to
    save yourself. Roll two dice. If the score is greater than your
    Rank, you are drowned.

    [Next page](680)
    """
    else:
        page["markdown"] = """
    Your ship, crew and cargo are lost to the deep, dark sea. Cross
    them off your Adventure Sheet. Your only thought now is to
    save yourself. 

    You manage to find some driftwood, and
    make it back to shore. Lose %d Stamina point(s).

    [Next page](173)
    """
