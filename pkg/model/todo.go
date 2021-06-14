package model

import "time"

type Todo struct {
	Id          *string    `json:"id" validate:"required"`
	Name        *string    `json:"name" validate:"required"`
	CreatedTime *time.Time `json:"createdTime" validate:"required"`
	UpdatedTime *time.Time `json:"updatedTime"`
	Completed   bool       `json:"completed"`
}
