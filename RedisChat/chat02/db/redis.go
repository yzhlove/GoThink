package db

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var pool *redis.Pool

// Init ...
func Init() error {
	pool = &redis.Pool{
		MaxIdle:     10 * runtime.NumCPU(),
		MaxActive:   50 * runtime.NumCPU(),
		Wait:        true,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				"127.0.0.1:6379",
				redis.DialDatabase(0),
				redis.DialConnectTimeout(5*time.Second),
				redis.DialReadTimeout(5*time.Second),
				redis.DialWriteTimeout(5*time.Second),
			)
		},
	}

	var conn redis.Conn

	for {
		conn = pool.Get()
		if err := conn.Err(); err != nil {
			fmt.Printf("redis connect failed with error [%v],cfg [%#v]", err, "127.0.0.1:6379")
			_ = conn.Close()
			time.Sleep(2 * time.Second)
			continue
		}
	RETRY:
		if _, err := conn.Do("PING"); err != nil {
			fmt.Println("wait for redis up")
			time.Sleep(2 * time.Second)
			goto RETRY
		} else {
			break
		}
	}

	return conn.Close()
}

func Key(mode string, uid uint64, args ...string) string {
	key := mode + ":{" + strconv.FormatUint(uid, 10) + "}"
	if len(args) > 0 {
		key += ":" + strings.Join(args, ":")
	}
	return key
}

func Pool() *redis.Pool {
	return pool
}

func Conn() redis.Conn {
	return pool.Get()
}
