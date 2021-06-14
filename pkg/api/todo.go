package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
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

	rd := c.Get("Redis").(*redis.Client)

	val, err := rd.Get(ctx, id).Result()

	if err != nil {
		if err == redis.Nil {
			return echo.NewHTTPError(http.StatusNotFound, "No todo item found")
		} else {
			panic(err)
		}
	}

	// t := model.Todo{
	// 	Id:          new(string),
	// 	CreatedTime: new(time.Time),
	// }

	// *t.Id = uuid.New().String()
	// *t.CreatedTime = time.Now().UTC()

	var t model.Todo

	json.Unmarshal([]byte(val), &t)

	return c.JSON(http.StatusOK, t)
}

func UpdateTodo(c echo.Context) error {
	id := c.Param("id")

	rd := c.Get("Redis").(*redis.Client)

	val, err := rd.Get(ctx, id).Result()

	if err != nil {
		if err == redis.Nil {
			return echo.NewHTTPError(http.StatusNotFound, "No todo item found")
		} else {
			panic(err)
		}
	}

	var existingTodo model.Todo

	json.Unmarshal([]byte(val), &existingTodo)

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

	b, err := json.Marshal(existingTodo)

	if err != nil {
		panic(err.Error())
	}

	err = rd.Set(ctx, *existingTodo.Id, b, time.Minute*5).Err() //Cache data for 5 minutes

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

	b, err := json.Marshal(*t)

	if err != nil {
		panic(err.Error())
	}

	rd := c.Get("Redis").(*redis.Client)

	err = rd.Set(ctx, *t.Id, b, time.Minute*5).Err() //Cache data for 5 minutes

	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, t)
}
