package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/kabelsea-sandbox/kun"
	"github.com/kabelsea-sandbox/kun/model"
)

// FindOp interface for repository
type FindOp[T model.Model] interface {
	Operator
	Find(ctx context.Context, filters ...FilterOperator) (T, error)
}

// find implement FindOp interface
type find[T model.Model] struct {
	client *bun.DB
}

// NewFindOp returns new find operator
func NewFindOp[T model.Model](client *bun.DB) FindOp[T] {
	return &find[T]{
		client: client,
	}
}

// Find implement FindOp interface
func (op *find[T]) Find(ctx context.Context, filters ...FilterOperator) (instance T, err error) {
	operator := &findOperator[T]{
		operatorContext: op,
	}
	err = operator.Execute(ctx, instance, filters...)
	return
}

// findOperator type provide API for execute query
type findOperator[T model.Model] OperatorContext[*find[T]]

func (op *findOperator[T]) Execute(ctx context.Context, instance T, filters ...FilterOperator) error {
	query := op.operatorContext.client.
		NewSelect().
		Model(instance)

	_, err := FilterMerge(query, filters...).Exec(ctx)

	return kun.HandleError(err)
}
