package api

import (
	"context"
	"net/http"
	"time"
	"todo/pkg/db"
	"todo/pkg/model"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var ctx = context.Background()

func GetTodos(c echo.Context) error {
	return c.JSON(http.StatusOK, "Todos")
}

func GetTodoById(c echo.Context) error {
	id := c.Param("id")

	d := c.Get("Db").(db.Database)

	t, err := d.GetTodo(id)

	if err != nil {
		if err == redis.Nil {
			return echo.NewHTTPError(http.StatusNotFound, "No todo item found")
		} else {
			panic(err)
		}
	}

	return c.JSON(http.StatusOK, t)
}

func UpdateTodo(c echo.Context) error {
	id := c.Param("id")

	d := c.Get("Db").(db.Database)

	existingTodo, err := d.GetTodo(id)

	if err != nil {
		if err == redis.Nil {
			return echo.NewHTTPError(http.StatusNotFound, "No todo item found")
		} else {
			panic(err)
		}
	}

	updatedTodo := new(model.Todo)

	if err := c.Bind(&updatedTodo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingTodo.Name = updatedTodo.Name
	existingTodo.CreatedTime = updatedTodo.CreatedTime

	if existingTodo.UpdatedTime == nil {
		existingTodo.UpdatedTime = new(time.Time)
	}

	*existingTodo.UpdatedTime = time.Now().UTC()
	existingTodo.Completed = updatedTodo.Completed

	if err := c.Validate(existingTodo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = d.AddTodo(id, existingTodo)

	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, existingTodo)
}

func AddTodo(c echo.Context) error {

	t := new(model.Todo)

	if err := c.Bind(&t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	t.Id = new(string)
	*t.Id = uuid.New().String()

	t.CreatedTime = new(time.Time)
	*t.CreatedTime = time.Now().UTC()

	if err := c.Validate(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	d := c.Get("Db").(db.Database)

	err := d.AddTodo(*t.Id, *t)

	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, t)
}
