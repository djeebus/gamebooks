load("../lib/stamina.star", "stamina_add")

markdown = """
Someone stabs you in the back. Lose 5 Stamina points. If you
still live, you spin around just as a beefy, disreputable-looking
thug comes for you again with a long dagger.

‘Get the snooping swine!’ yells the man with the eyepatch.

You must fight.

[Thug, COMBAT 3, Defence 6, Stamina 13](!fight!3!6!13!476)
"""

def once():
    stamina_add(-5)
