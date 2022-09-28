package redisHelper

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func SaveKey(context context.Context, key string, value interface{}, expire time.Duration) *redis.StatusCmd {
	rdb := GetRedis()
	redisDB := rdb.GetConnection()
	defer rdb.Release(redisDB)

	return redisDB.Set(context, key, value, expire)

}

func SaveOnTable(context context.Context, key string, value ...interface{}) *redis.IntCmd {
	rdb := GetRedis()
	redisDB := rdb.GetConnection()
	defer rdb.Release(redisDB)

	return redisDB.HSet(context, key, value)

}

func GetFromTable(context context.Context, key string, field string) *redis.StringCmd {
	rdb := GetRedis()
	redisDB := rdb.GetConnection()
	defer rdb.Release(redisDB)

	return redisDB.HGet(context, key, field)
}

func GetValue(context context.Context, key string) *redis.StringCmd {
	rdb := GetRedis()
	redisDB := rdb.GetConnection()
	defer rdb.Release(redisDB)

	return redisDB.Get(context, key)
}

func GetTTL(context context.Context, key string) *redis.DurationCmd {
	rdb := GetRedis()
	redisDB := rdb.GetConnection()
	defer rdb.Release(redisDB)

	return redisDB.TTL(context, key)
}

func GetAll(context context.Context, key string) *redis.StringStringMapCmd {
	rdb := GetRedis()
	redisDB := rdb.GetConnection()
	defer rdb.Release(redisDB)
	return redisDB.HGetAll(context, key)
}

func DeleteKeys(ctx context.Context, key ...string) *redis.IntCmd {
	rdb := GetRedis()
	redisDB := rdb.GetConnection()
	defer rdb.Release(redisDB)
	return redisDB.Del(ctx, key...)
}

func DeleteFromTable(ctx context.Context, key string, field ...string) *redis.IntCmd {
	rdb := GetRedis()
	redisDB := rdb.GetConnection()
	defer rdb.Release(redisDB)
	return redisDB.HDel(ctx, key, field...)
}
