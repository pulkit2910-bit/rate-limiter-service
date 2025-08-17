local userId = KEYS[1]
local apiPath = KEYS[2]
local currentTokensArg = ARGV[1]
local lastUpdatedArg = ARGV[2]
local capacityArg = ARGV[3]
local refillRateArg = ARGV[4]
local nowTime = ARGV[5]

-- Config values are stored in a separate key
local configKey = "config:" .. userId
local configValue = redis.call("HMGET", configKey, capacityArg, refillRateArg)
if configValue[1] == false or configValue[2] == false then
    return { 0, "Configuration not found for userId: " .. userId }
end

local allowed = 0
local tokens = tonumber(configValue[1])
local lastUpdatedTime = tonumber(nowTime)

local key = userId .. ":" .. apiPath
local currentValue = redis.call("HMGET", key, currentTokensArg, lastUpdatedArg)

-- If the current value is nil, it means no tokens have been set yet
if currentValue[1] == false or currentValue[2] == false then
    allowed = 1
    redis.call("HSET", key, currentTokensArg, tokens, lastUpdatedArg, lastUpdatedTime)
    return { allowed, tokens, lastUpdatedTime }
end

local capactity = tonumber(configValue[1])
local refillRate = tonumber(configValue[2])

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
redis.call("HSET", key, currentTokensArg, tokens, lastUpdatedArg, lastUpdatedTime)

return { allowed, tokens, lastUpdatedTime }