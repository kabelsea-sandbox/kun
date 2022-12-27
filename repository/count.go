package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/kabelsea-sandbox/kun"
)

// CountOp public interface
type CountOp[T any] interface {
	Count(ctx context.Context, filters ...FilterOperator) (int, error)
}

// CountOp private interface
type countOp[T any, PT ptr[T]] interface {
	Operator
	CountOp[T]
}

// Count implement CountOp interface
type count[T any, PT ptr[T]] struct {
	client *bun.DB
}

func NewCountOp[PT ptr[T], T any](client *bun.DB) CountOp[PT] {
	return &count[T, PT]{
		client: client,
	}
}

// Count implement CountOp interface
func (o *count[T, PT]) Count(ctx context.Context, filters ...FilterOperator) (int, error) {
	operator := &countOperator[T, PT]{
		operatorContext: o,
	}

	return operator.Execute(ctx, filters...)
}

type countOperator[T any, PT ptr[T]] OperatorContext[*count[T, PT]]

func (op *countOperator[T, PT]) Execute(ctx context.Context, filters ...FilterOperator) (int, error) {
	query := op.operatorContext.client.
		NewSelect().
		Model(new(T))

	res, err := FilterMerge(query, filters...).Count(ctx)

	return res, kun.HandleError(err)
}
