package middleware

import (
	"todo/pkg/db"

	"github.com/labstack/echo/v4"
)

const (
	TxKey = "Db"
)

func RedisHandler(d db.Database) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.Set(TxKey, d)

			return next(c)
		})
	}
}
