package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/kabelsea-sandbox/kun"
	"github.com/kabelsea-sandbox/kun/model"
)

const (
	searchOpDefaultLimit = 30
)

// SearchOp interface for repository
type SearchOp[T model.Model] interface {
	Operator
	Search(ctx context.Context, page Pager, filters ...FilterOperator) ([]T, *PageInfo, error)
}

// Search implement SearchOp interface
type search[T model.Model] struct {
	client *bun.DB
}

func NewSearchOp[T model.Model](client *bun.DB) SearchOp[T] {
	return &search[T]{
		client: client,
	}
}

// Search implement SearchOp interface
func (op *search[T]) Search(
	ctx context.Context,
	page Pager,
	filters ...FilterOperator,
) ([]T, *PageInfo, error) {
	operator := &searchOperator[T]{
		operatorContext: op,
	}
	return operator.Execute(ctx, page, filters...)
}

type searchOperator[T model.Model] OperatorContext[*search[T]]

// Execute query
func (op *searchOperator[T]) Execute(
	ctx context.Context,
	page Pager,
	filters ...FilterOperator,
) ([]T, *PageInfo, error) {
	var (
		result = make([]T, 0, searchOpDefaultLimit)
	)

	query := op.operatorContext.client.
		NewSelect().
		Model(&result)

	if page != nil {
		filters = append(filters, Page(page))
	}

	// filters
	query = FilterMerge(query, filters...)

	var (
		count int
		err   error
	)

	if page != nil && !page.GetSkipTotal() {
		count, err = query.ScanAndCount(ctx)
	} else {
		err = query.Scan(ctx)
	}

	if err != nil {
		return nil, nil, kun.HandleError(err)
	}

	start := StartCursor(result)
	end := EndCursor(result)

	var pageInfo = &PageInfo{
		Start:       start,
		End:         end,
		Total:       count,
		HasNext:     !(end == ""),
		HasPrevious: !(start == ""),
	}

	return result, pageInfo, nil
}
