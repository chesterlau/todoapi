package db

import (
	"context"
	"encoding/json"
	"time"
	"todo/pkg/model"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Address string
}

var ctx = context.Background()
var client *redis.Client

func (r Redis) Init() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	client = rdb
}

func (r Redis) AddTodo(id string, t model.Todo) error {

	b, err := json.Marshal(t)

	if err != nil {
		panic(err.Error())
	}

	err = client.Set(ctx, id, b, time.Minute*5).Err() //Cache data for 5 minutes

	if err != nil {
		return err
	}

	return nil
}

func (r Redis) GetTodo(id string) (model.Todo, error) {

	val, err := client.Get(ctx, id).Result()

	var t model.Todo

	if err != nil {
		return t, err
	}

	json.Unmarshal([]byte(val), &t)

	return t, nil
}

func (r Redis) GetTodos() ([]model.Todo, error) {

	val, err := client.Keys(ctx, "*").Result()

	ts := make([]model.Todo, 0)

	if err != nil {
		return ts, err
	}

	for _, key := range val {

		var t model.Todo

		v, err := client.Get(ctx, key).Result()

		if err != nil {
			return ts, err
		}

		json.Unmarshal([]byte(v), &t)

		ts = append(ts, t)
	}

	return ts, nil
}

func (r Redis) DeleteTodo(id string) (int64, error) {

	val, err := client.Del(ctx, id).Result()

	if err != nil {
		return 0, err
	}

	return val, nil
}
