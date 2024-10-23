local abilities = require("lib/abilities")
local bank = require("lib/bank")

local _M = {
    name = "Star smuggler",
    start_page = "e001",

    on_start = function()
        abilities.init()
        bank.init()
    end,

    on_page = function(page)
        header = [[

| attributes | value |
| ---------- | ----- |
| cunning | ]] .. abilities.get_cunning() .. [[ |
| cash | ]] .. bank.get_cash() .. [[ |

]]

        page.markdown = header .. page.markdown

        return page
    end,
}

return _M
