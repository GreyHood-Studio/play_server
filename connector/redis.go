package connector

import (
	"github.com/go-redis/redis"
	"fmt"
)

type RedisConn struct {
	client	*redis.Client
	addr	string
	pwd		string
	db		int
}

// Redis Connect 로직
func (conn *RedisConn) connectRedis() {
	conn.client = redis.NewClient(&redis.Options{
		Addr:     conn.addr,
		Password: conn.pwd, // no password set
		DB:       conn.db,  // use default DB
	})

	pong, err := conn.client.Ping().Result()
	// 나중에 logger로 바꿀 것
	fmt.Println(pong, err)
	// Output: PONG <nil>
}