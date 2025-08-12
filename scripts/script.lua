local ttl = redis.Call("GET", KEYS[1])

if ttl < 30 then
    return "soon"
elseif ttl >= 30 then
    return "ok"
elseif ttl == -1 then
    return "no expiry"
else 
    return "no key found"
end