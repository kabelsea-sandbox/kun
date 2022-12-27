package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/kabelsea-sandbox/kun"
	"github.com/kabelsea-sandbox/kun/model"
)

// DeleteOp public interface
type DeleteOp[T any] interface {
	Delete(ctx context.Context, id model.ID) error
}

// DeleteOp private interface
type deleteOp[T any, PT ptr[T]] interface {
	Operator
	DeleteOp[T]
}

// Delete implement DeleteOp interface
type delete[T any, PT ptr[T]] struct {
	client *bun.DB
}

func NewDeleteOp[PT ptr[T], T any](client *bun.DB) DeleteOp[PT] {
	return &delete[T, PT]{
		client: client,
	}
}

// Delete implement DeleteOp interface
func (o *delete[T, PT]) Delete(ctx context.Context, id model.ID) error {
	operator := &deleteOperator[T, PT]{
		operatorContext: o,
	}

	return operator.Execute(ctx, id)
}

type deleteOperator[T any, PT ptr[T]] OperatorContext[*delete[T, PT]]

func (op *deleteOperator[T, PT]) Execute(ctx context.Context, id model.ID) error {
	_, err := op.operatorContext.client.
		NewDelete().
		Model(new(T)).
		Where("id = ?", id).
		Exec(ctx)

	return kun.HandleError(err)
}
