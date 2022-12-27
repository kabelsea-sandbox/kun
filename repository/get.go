package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/kabelsea-sandbox/kun"
	"github.com/kabelsea-sandbox/kun/model"
)

// GetOp public interface
type GetOp[T model.Model] interface {
	Get(ctx context.Context, id model.ID) (T, error)
}

// GetOp private interface
type getOp[T any, PT ptr[T]] interface {
	Operator
	GetOp[PT]
}

// get implement Getter interface
type get[T any, PT ptr[T]] struct {
	client *bun.DB
}

// NewGetOp returns new get operator
func NewGetOp[PT ptr[T], T any](client *bun.DB) GetOp[PT] {
	return &get[T, PT]{
		client: client,
	}
}

// Get implement Getter interface
func (op *get[T, PT]) Get(ctx context.Context, id model.ID) (PT, error) {
	operator := &getOperator[T, PT]{
		operatorContext: op,
	}
	return operator.Execute(ctx, id)
}

// getOperator type provide API for execute query
type getOperator[T any, PT ptr[T]] OperatorContext[*get[T, PT]]

func (op *getOperator[T, PT]) Execute(ctx context.Context, id model.ID) (PT, error) {
	instance := PT(new(T))

	err := op.operatorContext.client.
		NewSelect().
		Model(instance).
		Where("id = ?", id).
		Scan(ctx)

	return instance, kun.HandleError(err)
}
