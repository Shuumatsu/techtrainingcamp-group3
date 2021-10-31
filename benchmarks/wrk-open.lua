local json = require 'json'

local path = "/open"
local method = "POST"

request = function()
    local uid = math.randomseed(os.time())
    local envelope_id = math.randomseed(os.time())

    local headers = {
        ["Content-Type"] = "application/json",
        ["Accept"] = "application/json"
    }

    local body = json.encode({uid = uid, envelope_id = envelope_id})

    return wrk.format(method, path, headers, body)
end

