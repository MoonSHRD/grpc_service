package dbs

import (
    "github.com/go-redis/redis"
)


type redi struct {
    instance *redis.Client
}

var redisInstance *redi

func GetRedis() *redis.Client {
    if mongoInstance == nil {
        conf:=getConfig("redis")
        client := redis.NewClient(&redis.Options{
            Addr:     conf.Host+":"+conf.Port,
            Password: conf.Password,
            DB:       0,
        })
        redisInstance = &redi{instance: client}
    }
    return redisInstance.instance
}