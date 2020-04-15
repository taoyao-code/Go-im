package server

//
//import (
//	"log"
//	"net/http"
//	"reptile-go/util"
//	"time"
//
//	"github.com/go-redis/redis"
//
//	"github.com/gin-gonic/gin"
//)
//
//var Client *redis.Client
//
//func init() {
//	Client = redis.NewClient(&redis.Options{
//		Addr:         "127.0.0.1:6379",
//		PoolSize:     1000, // 池子
//		ReadTimeout:  time.Millisecond * time.Duration(100),
//		WriteTimeout: time.Millisecond * time.Duration(100),
//		IdleTimeout:  time.Millisecond * time.Duration(60), // 空闲超时
//	})
//	_, err := Client.Ping().Result()
//	if err != nil {
//		panic("init redis error")
//	}
//}
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
