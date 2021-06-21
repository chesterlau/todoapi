package router

import (
	"net/http"
	"todo/pkg/api"
	"todo/pkg/db"
	"todo/pkg/middleware"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func Init() *echo.Echo {
	e := echo.New()

	//Custom middleware
	r := db.Redis{
		Address: "127.0.0.1:6379",
	}

	r.Init()

	e.Use(middleware.RedisHandler(r))

	v := validator.New()
	e.Validator = &CustomValidator{v}

	e.GET("/todos", api.GetTodos)
	e.GET("/todos/:id", api.GetTodoById)
	e.POST("/todos", api.AddTodo)
	e.PUT("/todos/:id", api.UpdateTodo)
	e.DELETE("/todos/:id", api.DeleteTodo)

	return e
}
