local json = require 'json'

local path = '/snatch'
local method = 'POST'

request = function()
    local uid = math.randomseed(os.time())

    local headers = {
        ["Content-Type"] = "application/json",
        ["Accept"] = "application/json"
    }

    local body = json.encode({uid = uid})

    return wrk.format(method, path, headers, body)
end

