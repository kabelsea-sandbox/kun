package repository

import "github.com/kabelsea-sandbox/kun/model"

type ptr[T any] interface {
	*T
	model.Model
}
