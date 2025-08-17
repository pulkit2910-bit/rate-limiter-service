local userId = KEYS[1]
local apiPath = KEYS[2]
local currentTokensArg = ARGV[1]
local lastUpdatedArg = ARGV[2]
local capactityArg = ARGV[3]
local refillRateArg = ARGV[4]
local nowTime = ARGV[5]

local key = userId .. ":" .. apiPath
local currentValue = redis.Call("HGET", key, currentTokensArg, lastUpdatedArg)

-- Config values are stored in a separate key
local configKey = "config:" .. userId
local configValue = redis.Call("HGET", configKey, capacityArg, refillRateArg)

local capactity = tonumber(configValue[1])
local refillRate = tonumber(congigValue[2])

local allowed = 0
local tokens = tonumber(configValue[1])
local lastUpdatedTime = tonumber(nowTime)

-- If the current value is nil, it means no tokens have been set yet
if currentValue == false then
    allowed = 1
    redis.Call("HSET", key, currentTokensArg, tokens, lastUpdatedArg, lastUpdatedTime)
    return { allowed, tokens, lastUpdatedTime }
end

tokens = tonumber(currentValue[1])
lastUpdatedTime = tonumber(currentValue[2])

-- Calculate the elapsed time since the last update and refill the tokens
local elapsed = nowTime - lastUpdatedTime
local newTokens = math.min(tokens + (elapsed / refillRate), capactity)

-- If the new tokens are greater than 0, allow the request and decrement the token count
if newTokens > 0 then
    allowed = 1
    newTokens = newTokens - 1
else
    allowed = 0
end

tokens = newTokens
lastUpdatedTime = tonumber(nowTime)
redis.Call("HSET", key, currentTokensArg, tokens, lastUpdatedArg, lastUpdatedTime)

return { allowed, tokens, lastUpdatedTime }