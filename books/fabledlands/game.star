load("./lib/abilities.star", "ability_get")
load("./lib/bank.star", bank_balance="balance")
load("./lib/character.star", "defense_get", "name_get", "profession_get", "rank_get")
load("./lib/codewords.star", "codeword_assert", "codeword_all")
load("./lib/god.star", "god_get")
load("./lib/stamina.star", stamina_get="get", stamina_get_max="get_max")

start_page = "0"

def _require_codeword(codeword, page_id):
    codeword_assert(codeword)
    return page_id


def _make_roll(ability, difficulty, success_page_id, fail_page_id):
    dice = dice_roll(2, 6)
    ability = ability_get(ability)
    if dice + ability > int(difficulty):
        return success_page_id
    else:
        return fail_page_id


_base_commands = {
    'require-codeword': _require_codeword,
    'roll': _make_roll,
}

def on_command(command):
    parts = command.split('!')
    cmd = parts[0]
    args = parts[1:]

    fn = _base_commands.get(cmd)
    if fn:
        return fn(*args)

    fail(cmd, args)


def _codewords():
    result = ""
    col_count = 4
    codewords = codeword_all()
    for index, (codeword, value) in enumerate(codewords.items()):
        if index % col_count == 0:
            result += "<tr>"
        result += "<td><input type='checkbox' disabled %s />%s</td>" % (
            "checked" if value else "", codeword,
        )
        if index % col_count == col_count-1:
            result += "</tr>"

    return result


def on_page(page):
    if "allow_return" not in page:
        page["allow_return"] = False
    if "clear_history" not in page:
        page["clear_history"] = True
    if "on_command" not in page:
        page["on_command"] = on_command

    if page["page_id"] != start_page:
        markdown = page["markdown"]

        header = """
<table>
    <tbody>
        <tr>
            <th>Name</th>
            <th colspan="2">Profession</th>
        </tr>
        <tr>
            <td>%s</td>
            <td colspan="2">%s</td>
        </tr>
        <tr>
            <th>God</th>
            <th>Rank</th>
            <th>Defence</th>
        </tr>
        <tr>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
        </tr>
        <tr>
            <td>
                <table>
                    <thead>
                        <tr>
                            <th>Ability</th>
                            <th>Score</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>Charisma</td>
                            <td>%s</td>
                        </tr>
                        <tr>
                            <td>Combat</td>
                            <td>%s</td>
                        </tr>
                        <tr>
                            <td>Magic</td>
                            <td>%s</td>
                        </tr>
                        <tr>
                            <td>Sanctity</td>
                            <td>%s</td>
                        </tr>
                        <tr>
                            <td>Scouting</td>
                            <td>%s</td>
                        </tr>
                        <tr>
                            <td>Thievery</td>
                            <td>%s</td>
                        </tr>
                    </tbody>
                </table>
                <table>
                    <thead>
                        <tr>
                            <th colspan="2">Stamina</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>When unwounded</td>
                            <td>%s</td>
                        </tr>
                        <tr>
                            <td>Current:</td>
                            <td>%s</td>
                        </tr>
                    </tbody>
                </table>
            </td>
            <td>
                <table>
                    <thead>
                        <tr>
                            <th>Possessions (maximum of 12)</th>
                        </tr>
                    </thead>
                    <tbody>%s</tbody>
                </table>
                <table>
                    <thead>
                        <tr>
                            <th>Codewords</th>
                        </tr>
                    </thead>
                    <tbody>%s</tbody>
                </table>
            </td>
        </tr>
        <tr>
            <td>
                <table> 
                    <thead>
                        <tr>
                            <th>Resurrection Arrangements</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>%s</td>
                        </tr>
                    </tbody>
                </table>
            </td>
            <td>
                <table> 
                    <thead>
                        <tr>
                            <th>Money</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>%d</td>
                        </tr>
                    </tbody>
                </table>
            </td>
        </tr>
        <tr>
            <td>
                <table> 
                    <thead>
                        <tr>
                            <th>Titles and Honours</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>%s</td>
                        </tr>
                    </tbody>
                </table>
            </td>
            <td>
                <table> 
                    <thead>
                        <tr>
                            <th>Blessings</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>%s</td>
                        </tr>
                    </tbody>
                </table>
            </td>
        </tr>
    </tbody>
</table>
    """ % (
            name_get(),
            profession_get(),
            god_get(),
            rank_get(),
            defense_get(),
            ability_get('charisma'),
            ability_get('combat'),
            ability_get('magic'),
            ability_get('sanctity'),
            ability_get('scouting'),
            ability_get('thievery'),
            stamina_get_max(),
            stamina_get(),
            "", # possessions
            _codewords(), # codewords
            "", # resurrection
            bank_balance(),
            "", # title and honours
            "", # blessings
        )

        page["markdown"] = "%s\n\n%s" % (header, markdown)
