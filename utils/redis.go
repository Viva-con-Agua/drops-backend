package utils

import (
	"fmt"
	"github.com/Viva-con-Agua/echo-pool/pool"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/rbcervilla/redisstore"
	"log"
	"net/http"
)

func RedisSession() echo.MiddlewareFunc {
	client := redis.NewClient(&redis.Options{
		Addr: "172.2.150.2:6379",
	})

	redis, err := redisstore.NewRedisStore(client)

	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}
	fmt.Println("Redis successfully connected!")
	return session.Middleware(redis)
}

func SessionAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		if sess.Values["valid"] == nil || sess.Values["valid"] == false {
			return echo.NewHTTPError(http.StatusUnauthorized, pool.Unauthorized())
		}
		return next(c)
	}
}
