package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/kabelsea-sandbox/kun"
	"github.com/kabelsea-sandbox/kun/model"
)

// ExistsOp interface for repository
type ExistsOp[T model.Model] interface {
	Operator
	Exists(ctx context.Context, filters ...FilterOperator) (bool, error)
}

// Exists implement ExistsOp interface
type exists[T model.Model] struct {
	client *bun.DB
}

func NewExistsOp[T model.Model](client *bun.DB) ExistsOp[T] {
	return &exists[T]{
		client: client,
	}
}

// Exists implement ExistsOp interface
func (o *exists[T]) Exists(ctx context.Context, filters ...FilterOperator) (bool, error) {
	operator := &existsOperator[T]{
		operatorContext: o,
	}

	return operator.Execute(ctx, filters...)
}

type existsOperator[T model.Model] OperatorContext[*exists[T]]

func (op *existsOperator[T]) Execute(ctx context.Context, filters ...FilterOperator) (bool, error) {
	query := op.operatorContext.client.
		NewSelect().
		Model((*T)(nil))

	exists, err := FilterMerge(query, filters...).Exists(ctx)

	return exists, kun.HandleError(err)
}
