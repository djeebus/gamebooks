load("../lib/codewords.star", "codeword_add", "codeword_has")

markdown = """
If you have the codeword Altitude, turn to 272 immediately. If
not, read on.

A notice has been pinned up in the foyer. ‘Adventurer priest
wanted. See the Chief Administrator.’

Naturally, you present yourself, and the Chief Administrator,
a grey-whiskered priest of Elnir, takes you into his office. He
shows you a special crystal ball that displays an aerial view of
Marlock City. You notice several strange-looking clouds
hanging over the city. They are shaped like gigantic demons
reaching down to claw at the city laid out below them.

‘The crystal ball shows things as they are in the spirit world,’
explains the priest. ‘These storm demons cannot be seen under
normal circumstances, but they are there, almost ready to
destroy the city.’

He goes on to tell you that Sul Veneris, the divine Lord of
Thunder is one of the sons of Elnir, the Sky God, chief among
the gods. He is responsible for keeping the storm demons under
control, and thunder is thought to be the sound of Sul Veneris
smiting the demons in his wrath.

‘Unfortunately, the storm demons have found a way to put
Sul Veneris into an enchanted sleep. He lies at the very top of
Devil’s Peak, a single spire of volcanic rock, reaching up into the
clouds. The peak lies north of Marlock City and the Curstmoor.

We need an enterprising priest to get to the top of the peak and
free Sul Veneris from his sleep. But I must warn you that several
priests have already tried, and we never saw them again.’

If you take up the quest, [record the codeword Altitude](!accept). Otherwise [continue](100).
"""

def on_page(page):
    if codeword_has('altitude'):
        return "272"

def on_command(command):
    if command == 'accept':
        codeword_add("altitude")
        return "100"
