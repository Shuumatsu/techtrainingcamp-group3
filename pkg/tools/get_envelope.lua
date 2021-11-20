local TotalAmount = redis.call('GET', 'TotalAmount')
local EnvelopeAmount = redis.call('INCR','EnvelopeAmount')

if (EnvelopeAmount > TotalAmount)
then
    return 0
end
local LeftMoney = redis.call('GET', 'TotalMoney') - redis.call('GET', 'UsedMoney')

if (LeftMoney <= 0)
then
    return 0
end

local MinMoney = redis.call('GET', 'MinMoney')

if (MinMoney > LeftMoney)
then
    return 0
end

local MaxMoney = math.min(redis.call('GET', 'MaxMoney'), LeftMoney)

math.randomseed(os.time())
local Money =  math.random(MinMoney, MaxMoney)

redis.call('INCRBY', 'UsedMoney', Money)
redis.call('INCR', 'EnvelopeAmount')

return Money


