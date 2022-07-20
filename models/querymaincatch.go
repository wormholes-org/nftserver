package models

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

var mqr *queryRedisCatch

func NewQueryMainCatch(redisIP, passwd string) error {
	mqr = new(queryRedisCatch)
	mqr.pool = &redis.Pool{
		MaxActive:   queryRedisMaxactive,
		MaxIdle:     queryRedisMaxidle,
		IdleTimeout: queryRedisIdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisIP, redis.DialPassword(passwd))
		},
	}
	mqr.CathFlag = true
	mqr.connPool = make(map[string]Remote)
	mqr.close = make(chan struct{})
	return nil
}

func GetRedisMainCatch() *queryRedisCatch {
	if qr == nil {
		log.Fatalln("mqr is nil.")
	}
	return mqr
}
