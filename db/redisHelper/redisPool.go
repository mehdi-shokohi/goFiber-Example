package redisHelper

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	conf "goex/config"
)

var redisConnection *Redis

// Redis ...
type Redis struct {
	ctx       context.Context
	redisPool chan *redis.Client
}

// GetRedis
func (pool *Redis) GetConnection() *redis.Client {
	if pool.redisPool == nil {
		pool.initialize()
	}
	if len(pool.redisPool) == 0 {
		client := redis.NewClient(&redis.Options{
			Addr:     conf.GetKeydbAddress(),
			Password: conf.RedisPassword, // no password set
			DB:       0,                  // use default DB
		})
		pool.redisPool <- client
	}
	con := <-pool.redisPool
	_, err := con.Ping(pool.ctx).Result()
	if err != nil {
		panic(err)
	}

	return con

}

// GetRedis ...
func GetRedis() *Redis {
	if redisConnection == nil {
		fmt.Println("Redis RDB ....")
		redisConnection = new(Redis)
	}
	return redisConnection
}

// Release ...
func (pool *Redis) Release(con *redis.Client) {
	if len(pool.redisPool) > 500 {
		_ = con.Close()
	} else {
		pool.redisPool <- con
	}
}

func (pool *Redis) initialize() {
	fmt.Println("Redis Pool Initialized")
	pool.ctx = context.Background()
	pool.redisPool = make(chan *redis.Client, 1000)
	for range [4]int{} {
		client := redis.NewClient(&redis.Options{
			Addr:     conf.GetKeydbAddress(),
			Password: conf.RedisPassword, // no password set
			DB:       0,                  // use default DB
			PoolSize: 2,
		})
		pool.redisPool <- client
	}
}
