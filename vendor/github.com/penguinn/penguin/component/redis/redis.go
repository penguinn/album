package redis

import (
	"github.com/go-redis/redis"
	"log"
)

var (
	redisTemp *redis.Client
)

type ReidsConfig struct {
	Address        	string
	Password 		string
}

type RedisComponent struct {}

func (RedisComponent) Init(ops ...interface{}) (err error) {
	if len(ops) == 0 {
		errStr := "param error"
		log.Println("[redis error]:", errStr)
		return nil
	}
	cfg := ops[0].(*ReidsConfig)
	client := redis.NewClient(&redis.Options{
		Addr:cfg.Address,
		Password:cfg.Password,
	})
	redisTemp = client
	return nil
}

func Get() *redis.Client {
	return redisTemp
}