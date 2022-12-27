package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/kabelsea-sandbox/kun"
	"github.com/kabelsea-sandbox/kun/model"
)

// CreateOp interface for repository
type CreateOp[T model.Model] interface {
	Operator
	Create(ctx context.Context, create T) error
}

// Create implement CreateOp interface
type create[T model.Model] struct {
	client *bun.DB
}

// NewCreate returns create operator
func NewCreateOp[T model.Model](client *bun.DB) CreateOp[T] {
	return &create[T]{
		client: client,
	}
}

// Create implement CreateOp interface
func (o *create[T]) Create(ctx context.Context, create T) error {
	operator := &createOperator[T]{
		operatorContext: o,
	}

	return operator.Execute(ctx, create)
}

// createOperator type provide API for execute query
type createOperator[T model.Model] OperatorContext[*create[T]]

// Execute query end encode data to instance
func (op *createOperator[T]) Execute(ctx context.Context, instance T) error {
	_, err := op.operatorContext.client.
		NewInsert().
		Model(instance).
		Exec(ctx)

	return kun.HandleError(err)
}
