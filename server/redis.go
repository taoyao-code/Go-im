package server

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:         ViperConfig.Redis.Address,
		Password:     ViperConfig.Redis.Password,
		PoolSize:     1000,                                 // 池子
		ReadTimeout:  10 * time.Second,                     // 10秒
		WriteTimeout: 10 * time.Second,                     // 10秒
		IdleTimeout:  time.Millisecond * time.Duration(60), // 空闲超时
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis init data base ok")
}
func get(key string) (string, bool) {
	r, err := Client.Get(key).Result()
	if err == redis.Nil {
		fmt.Println(key + " does not exist")
	} else if err != nil {
		return "", false
	}
	return r, true
}
func set(key string, val string, expTime int32) {
	Client.Set(key, val, time.Duration(expTime)*time.Second)
}

//左边进
func Lpush(key string, val interface{}) {
	err := Client.LPush(key, val).Err()
	if err != nil {
		panic(err)
	}
}

//右边进
func Rpush(key string, val interface{}) {
	err := Client.RPush(key, val).Err()
	if err != nil {
		panic(err)
	}
}

// 左边出
func Lpop(key string) (string, bool) {
	r, err := Client.LPop(key).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		return "", false
	}
	return r, true
}

// 右边进右边出：栈
func Rpop(key string) (string, bool) {
	r, err := Client.RPop(key).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		return "", false
	}
	return r, true
}

// 从列表中获取数据
func LRange(key string, start, stop int64) ([]string, bool) {
	r, err := Client.LRange(key, start, stop).Result()
	if err == redis.Nil {
		return nil, false
	} else if err != nil {
		return nil, false
	}
	return r, true
}
