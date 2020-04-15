package server

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "foobared",
		PoolSize:     1000, // 池子
		ReadTimeout:  time.Millisecond * time.Duration(100),
		WriteTimeout: time.Millisecond * time.Duration(100),
		IdleTimeout:  time.Millisecond * time.Duration(60), // 空闲超时
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis init data base ok")
}

//func get(key string) (string, bool) {
//	r, err := Client.Get(key).Result()
//	if err != nil {
//		return "", false
//	}
//	return r, true
//}
//func set(key string, val string, expTime int32) {
//	Client.Set(key, val, time.Duration(expTime)*time.Second)
//}
//
///**
//set("name", "x", 100)
//s, b := get("name")
//fmt.Println(s, b)
//*/
//
////限制访问
//func RateMiddleware(limiter *util.Limiter) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 如果ip请求连接数在两秒内超过5次，返回429并抛出error
//		if !limiter.Allow(c.ClientIP(), 5, 2*time.Second) {
//			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
//				"error": "too many requests",
//			})
//			log.Println("too many requests")
//			return
//		}
//		c.Next()
//	}
//}
