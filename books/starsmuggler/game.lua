local abilities = require("lib/abilities")

local _M = {
    name = "Star smuggler",
    start_page = "e001",


    wrap_content = function(body)
        local content = [[
# Cool header content
]] .. body
        content = content .. [[
# Cool footer content
        ]]
        return content
    end,
}

function _M.on_start()
    abilities.init()
end

return _M
