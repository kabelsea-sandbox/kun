package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/kabelsea-sandbox/kun/model"
)

// BatchOp interface for repository
type BatchOp[T model.Model] interface {
	Operator
	Batch(ctx context.Context, ids ...model.ID) ([]T, error)
}

// Batcher implement batch interface
type batch[T model.Model] struct {
	client *bun.DB
}

// NewCreate returns create operator
func NewBatchOp[T model.Model](client *bun.DB) BatchOp[T] {
	return &batch[T]{
		client: client,
	}
}

// Batch implement Batcher interface
func (op *batch[T]) Batch(ctx context.Context, ids ...model.ID) (result []T, err error) {
	operator := &batchOperator[T]{
		operatorContext: op,
	}

	return operator.Execute(ctx, ids)
}

// batchOperator type provide API for execute query
type batchOperator[T model.Model] OperatorContext[*batch[T]]

// Execute query end encode data to instance
func (op *batchOperator[T]) Execute(ctx context.Context, ids []string) ([]T, error) {
	result := make([]T, 0, len(ids))

	err := op.operatorContext.client.
		NewSelect().
		Model(&result).
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)

	return result, err
}
