package middleware

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

const (
	TxKey = "Redis"
)

func RedisHandler(redis *redis.Client) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.Set(TxKey, redis)

			return next(c)
		})
	}
}
