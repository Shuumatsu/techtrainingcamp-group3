local json = require 'json'

local path = '/snatch'
local method = 'POST'

math.randomseed(os.time())

request = function()
    local uid = math.random(60000000)

    local headers = {
        ["Content-Type"] = "application/json",
        ["Accept"] = "application/json"
    }

    local body = json.encode({uid = uid})

    return wrk.format(method, path, headers, body)
end

