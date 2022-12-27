package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/kabelsea-sandbox/kun"
	"github.com/kabelsea-sandbox/kun/model"
)

// UpdateOp interface for repository
type UpdateOp[T model.Model] interface {
	Operator
	Update(ctx context.Context, update T) error
}

// Update implement Updater interface
type update[T model.Model] struct {
	client *bun.DB
}

func NewUpdateOp[T model.Model](client *bun.DB) UpdateOp[T] {
	return &update[T]{
		client: client,
	}
}

// Update implement Updater interface
func (op *update[T]) Update(ctx context.Context, update T) error {
	operator := &updateOperator[T]{
		operatorContext: op,
	}
	return operator.Execute(ctx, update)
}

type updateOperator[T model.Model] OperatorContext[*update[T]]

// Execute query
func (op *updateOperator[T]) Execute(ctx context.Context, update T) error {
	_, err := op.operatorContext.client.
		NewUpdate().
		Model(update).
		Where("id = ?", update.GetID()).
		Exec(ctx)

	return kun.HandleError(err)
}
