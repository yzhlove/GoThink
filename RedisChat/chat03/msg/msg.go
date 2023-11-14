package msg

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/reids-chat/chat03/db"
)

//go:generate msgp -tests=false -io=false

type Msg struct {
	Value   string
	Version int64
}

var redisAtomic = redis.NewScript(2, `
	local resp = redis.call('GET', KEYS[1])
if resp ~= false then 
    -- 尝试解码 MessagePack 数据，捕获任何错误
    local status, data = pcall(cmsgpack.unpack, resp)
    if not status then
        return redis.error_reply("Bad data format in input or other error: " .. data)
    end

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

func Set(key string, msg *Msg) error {

	bytes, err := msg.MarshalMsg(nil)
	if err != nil {
		return err
	}

	conn := db.Conn()
	defer conn.Close()

	_, err = redisAtomic.Do(conn, key, msg.Version, string(bytes))
	return err
}

func Get(key string) (*Msg, error) {

	conn := db.Conn()
	defer conn.Close()

	bytes, err := redis.String(conn.Do("GET", key))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return nil, nil
		}
		return nil, err
	}

	if len(bytes) != 0 {
		msg := &Msg{}
		if _, err := msg.UnmarshalMsg([]byte(bytes)); err != nil {
			return nil, err
		}
		return msg, nil
	}
	return nil, nil
}
