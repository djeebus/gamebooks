load("./lib/inventory.star", "inventory_add")
load("./lib/items.star", "all_items")

_potion_item_ids = [
    "truth-seeking-potion",
    "searching-potion",
    "transparency-potion",
    "invincibility-potion",
    "flying-potion",
    "power-potion",
    "calm-potion",
    "fear-potion",
    "intuition-potion",
    "duality-potion",
    "elusiveness-potion",
    "madness-potion",
    "revitalization-potion",
]

_cmd_select = "select"
_cmd_start = "start"
_potions_key = "potions"


def _select(potion_id):
    potion_ids = storage_page_get(_potions_key) or list()
    if potion_id in potion_ids:
        potion_ids = [pid for pid in potion_ids if pid != potion_id]
    else:
        potion_ids = list(potion_ids)
        potion_ids.extend([potion_id])
    storage_page_set(_potions_key, potion_ids)


def _start():
    potion_ids = storage_page_get(_potions_key) or list()
    if len(potion_ids) != 6:
        fail("must select 6 potions, only selected %d" % len(potion_ids))

    for potion_id in potion_ids:
        inventory_add(potion_id)

    return "pages/1"


def _render_table():
    selected_item_ids = storage_page_get(_potions_key) or list()

    rows = [
        "- [%s] [%s](!%s!%s): %s" % (
            "x" if item_id in selected_item_ids else " ",
            all_items[item_id]["name"],
            _cmd_select,
            item_id,
            all_items[item_id]["description"],
        ) for item_id in _potion_item_ids
    ]

    return "\n".join(rows)


markdown = """
In your fourth year at the Academy of the
Grand Wizard Eleutheria, you are becoming
bored. You have learned a great deal of magic,
of the power of reason, and the martial arts.
Now you yearn for a challenge. But Eleutheria
warns you that unless you stay for your fifth
and final year, you will not be fully proficient
in the arts of fighting, reason, and magic, and
will be vulnerable to evil forces.

However, you are adamant: you have spent
almost four years learning your skills, and are
acknowledged as the finest magician, best
thinker, and toughest fighter.

“This is true,” says Eleutheria. “But you still
have much to learn. Especially, you must
know how to decide whether reason, magic,
or force is the best course to take.”

But the prospect of another year in the
Academy is too much for you, and you decide
to leave. That night, you pack your book of
spells, enough food and drink for five meals,
and then look around the room you have lived
in for the past four years. It is more like a cell
than a room. You have had a stone bed to lie
on, and a Wizard’s dummy to practise your
sword-skills and Martial Arts against.

You reconsider: are you really going to leave
the security of the Academy? And where
would you go? You have heard stories of the
town of Darkblood, to the East, apparently
ruled over by a great but very evil Wizard.

You recall that Eleutheria had spoken briefly
only once about this Wizard, saying that he
was powerful and terrible, and almost the
equal of Eleutheria himself. For the Ruler of
Darkblood had once been one of the Grand
Wizards, but had been exiled when he used
his magic and power to plunder whole townships and steal treasure.

You go into the cell opposite, where your
friend Leofric is sitting on his bed, apparently
lost in thought. You know that he too has been
restless, and tell him that you are going to
leave that night in search of adventure. You
ask if he will accompany you.

At first he is enthusiastic, but then he has
second thoughts: “We would be alone, in
strange territory,” he says. “And as for Darkblood – there are rumours of adventurers
being buried alive or boiled in oil for even
daring to go there. We should at least wait
until Eleutheria has taught us all he can. One
day perhaps,” he says, “but not yet.”

You leave him in his cell with his dreams
and his fears, and walk back into your own
room. You will have to leave alone.

The idea in your mind becomes more exciting – to journey to Darkblood. You wonder
what you will find there. There are certainly
stories of the Ruler’s horrible vengeance
against those who have opposed him, and
you are far from sure that your own Powers
are enough yet to be a match for him. You are
strong; you have your Book of Spells – but you
know that, because of your inexperience, your
Spells do not always work.

However – the Magical potions kept in the
basement do always work in your experience.
You quickly decide that if you are to have any
chance against the Ruler of Darkblood, you
will need some Potions.

You creep down into the basement and look
at the bottles of Magic neatly arranged in an
alcove at the back. At first, you try to put as
many as you can into your backpack. But then
there is no room for your food and drink.
After packing and repacking several times,
you realise that you will have to be selective.
Only six Potions of Magic will fit into your
backpack. Choose any six different Potions.

%s

You choose six of the Potions. Each will
work only once, but you may drink a potion at
any time, except during a battle.

Wondering whether you have made the
right choices, you walk back up to your room,
and look round again.

Certainly you will need your sword, and
some kind of armour. You pick up your sword
and practise a few strokes. You feel confident.
You will go to the town of Darkblood and find
out whether the stories are true. Perhaps, you
may even meet the Ruler – and challenge him!

You buckle on your sword and shield and a
leather breastplate, and prepare to leave. Eleutheria’s words come back to you as you close
the door behind you: “. . . you must know
how to decide which is the best course to
take . . :’

For a moment, you hesitate; then resolute,
you stride confidently along the corridor,
down the stairs and out of the Academy. You
take the road East, wondering briefly whether
Eleutheria will use his Magic to transport you
back to the Academy. But nothing happens,
and you continue East until you reach the
town of Darkblood.

You walk through the gates into the town:
there seems to be a tremendous party in
progress. Creatures of all shapes and sizes are
swilling ale, laughing, singing, fighting, and
you stop in the main square to watch this
spectacle.

A strange, misshapen creature sidles up to
you. It is only about four feet high, and
although its left leg and arm are green and
warty, the toes and fingers webbed, the right
half of its body is perfect. The creature can
speak only with difficulty, the words slobbering from his mouth.

“Stranger,” he pleads, “we have not seen
such a warrior before in Darkblood for more
years than I can remember. You must take the
two jewels from the Wizard Belenghast, and
replace them where they rightly belong, in the
Temple of Valadon. Then will his hold be
released, and our town may once again return
to its former state of peace and tranquility.”

As the creature speaks, a figure lurches
towards you. It is cased in chain-mail, and
wears a helmet. Suddenly, the strange figure
swings a mace at you, but your four years of
training have more than prepared you for
such an attack.

You do not even draw your sword. Ducking
under the wild swing with the mace, you
strike a straight-fisted blow to the heart, and
as the figure falls to its knees, a double-fisted
blow behind the head.

The creature lies dead at your feet, and
again the half-man, half-toad appears at your
side.

“Come with me, Stranger,” he asks, and
you feel yourself unable to refuse. You follow
him to a dimly-lit but where sit a strange
collection of creatures. One has the body of a
snake, but the face of a man, another the body
of a crawling insect.

“I am Golfreth,” announces your companion. “Once I was a shield-maker, but now I
can no longer ply my trade. The same is true of
the rest of us. We have tried to enter the Tower
above, some of us have entered. Few have
returned, and those that have are in the pitiful
state you see. Until a warrior can be found
who can reach Belenghast, and replace the
two jewels in the Temple of Valadon, our city
will remain in chaos. Will you help us?”
You are excited and intrigued by the prospect, but you need more information. 

You learn that Belenghast is a great soldier who
also knows something of magic, and that he
lives in the tower beyond.

“You must climb to the second storey,”
intones the insect-man, “before you reach the
residence of Belenghast. And without the
Amulet of Stone, you will not succeed.”
You press the pitiful creatures for further
information, but there is little else they know.
Once you enter the Tower, there are many
traps and pitfalls. The route to Belenghast is
protected by strange creatures, some very
powerful, some cunning, and some skilled in
the use of Black Magic.

None has ever seen Belenghast, though he
is rumoured to be able to change his appearance at will. Only the fabled Amulet of Stone,
lost somewhere in the Tower, will reveal the
true Belenghast.

You leave the hut as the sun rises, and look
at the Tower ahead. It is a forbidding structure, and reminds you of a volcano. But you press on, determined, pushing your way
through the crazed townsfolk, until you reach
the foot of the Tower. A stone archway is
ahead of you, and beyond that is darkness.
You go back, push, your way into a Tavern,
and take a torch from the wall. Lighting it, you
walk back, through the archway, into the
Tower of Terror.

[Start the game](!start)
""" % (_render_table())


def on_command(cmd, args):
    if cmd == _cmd_select:
        return _select(*args)

    if cmd == _cmd_start:
        return _start()
