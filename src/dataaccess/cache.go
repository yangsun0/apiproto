package dataaccess

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheClient struct {
	redisClient *redis.Client
	context context.Context
}

func NewCacheClient() *CacheClient {
	return &CacheClient{}
}

func (cc *CacheClient) Connect(ctx context.Context) {
	opt := redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	}
	cc.redisClient = redis.NewClient(&opt)
	cc.context = ctx
}

func (cc *CacheClient) Close() {
	cc.redisClient.Close()
}

func (cc *CacheClient) Set(key string, value interface{}) {
	err := cc.redisClient.SetNX(cc.context, key, value, 10 * time.Second).Err()
	if err != nil {
		log.Printf("cache set error: %v", err)
	}
}

func (cc *CacheClient) Get(key string) string {
	value, err := cc.redisClient.Get(cc.context, key).Result()
	if err != nil {
		log.Printf("cache get error: %v", err)
		return ""
	}
	return value
}
