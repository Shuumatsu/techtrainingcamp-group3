local json = require 'json'

local path = "/open"
local method = "POST"

math.randomseed(os.time())

request = function()
    local uid = math.random(60000000)
    local envelope_id = math.random(2147483648)

    local headers = {
        ["Content-Type"] = "application/json",
        ["Accept"] = "application/json"
    }

    local body = json.encode({uid = uid, envelope_id = envelope_id})

    return wrk.format(method, path, headers, body)
end

