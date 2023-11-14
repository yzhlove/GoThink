package main

import "github.com/gomodule/redigo/redis"

/*
local resp = false

			if tonumber(ARGV[2]) > 0 then
				 resp = redis.call('SET', KEYS[1], ARGV[1],'NX', 'PX', ARGV[2])
			else
				 resp = redis.call('SET', KEYS[1], ARGV[1],'NX')
			end

			if resp == false then
				resp = redis.call('GET', KEYS[1])
				if resp ~= false then
					local uid = redis.call('HGET', KEYS[2],'PlayerId')
					if uid ~= false then
						local data = cmsgpack.unpack(resp)
						data['id'] = uid
						resp = cmsgpack.pack(data)
					end
				end
			end

			return resp
*/

var redisAtomic = redis.NewScript(2, `
	local resp = redis.call('GET', KEYS[1])
	if resp ~= false then 
    	local data = cmsgpack.unpack(resp)
    	local version = tonumber(data['Version'])
    	local expected_version = tonumber(KEYS[2])

   	 	if version and expected_version then
        	if version < expected_version then
            	return redis.call('SET', KEYS[1], ARGV[1])    
        	end
        	return false   
    	else
        	return redis.error_reply("Invalid 'version' or "..KEYS[2])
    	end
	else 
    	return redis.call('SET', KEYS[1], ARGV[1])
	end
`)
