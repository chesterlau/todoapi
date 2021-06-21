package db

import "todo/pkg/model"

type Database interface {
	Init()
	AddTodo(id string, t model.Todo) error
	GetTodo(id string) (model.Todo, error)
	DeleteTodo(id string) (int64, error)
}
