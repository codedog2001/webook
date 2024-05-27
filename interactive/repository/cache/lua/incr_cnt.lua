-- 具体业务
local key = KEYS[1]
-- 是阅读数，点赞数还是收藏数
local cntKey = ARGV[1]

local delta = tonumber(ARGV[2])

local exist=redis.call("EXISTS", key)
if exist == 1 then
    redis.call("HINCRBY", key, cntKey, delta)
    return 1
else
    return 0
end

--版本2： 当key不存在时，新建key，并将cntkey设置为1
---- 具体业务
--local key = KEYS[1]
---- 是阅读数，点赞数还是收藏数
--local cntKey = ARGV[1]
--
--local delta = tonumber(ARGV[2])
--
--local exist = redis.call("EXISTS", key)
--
--if exist == 1 then
--    redis.call("HINCRBY", key, cntKey, delta)
--else
--    -- 键不存在，新建键并设置cntKey字段为delta（此处为1）
--    redis.call("HSET", key, cntKey, delta)
--end
--
--return 1
