package repository

import (
	"github.com/uptrace/bun"
)

type Filtered interface{}

type FilteredCommon interface {
	Filtered
	~int | ~int32 | ~string | ~bool
}

type FilteredString interface {
	Filtered
	~string
}

type FilterOperator func(query *bun.SelectQuery) *bun.SelectQuery

type Filter[T FilteredCommon] func(field string, value T) FilterOperator

func Empty[T FilteredCommon]() FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query
	}
}

// IsNull filter
func IsNull(field string) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("? IS NULL", bun.Ident(field))
	}
}

// Equal filter
func Equal[T FilteredCommon](field string, value T) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("? = ?", bun.Ident(field), value)
	}
}

// NotEqual filter
func NotEqual[T FilteredCommon](field string, value T) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("? != ?", bun.Ident(field), value)
	}
}

// In filter
func In[T FilteredCommon](field string, values []T) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("? IN (?)", bun.Ident(field), bun.In(values))
	}
}

// InOr filter
func InOr[T FilteredCommon](field string, values []T) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.WhereOr("? IN (?)", bun.Ident(field), bun.In(values))
	}
}

// NotIn filter
func NotIn[T FilteredCommon](field string, values []T) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("? NOT IN (?)", bun.Ident(field), bun.In(values))
	}
}

// And group filter
func And(ops ...FilterOperator) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			for _, op := range ops {
				q = FilterMerge(q, op)
			}
			return q
		})
	}
}

// Or group filter
func Or(ops ...FilterOperator) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
			for _, op := range ops {
				q = FilterMerge(q, op)
			}
			return q
		})
	}
}

// Like filter
func Like[T FilteredString](field string, value T) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("? LIKE ?", bun.Ident(field), value+"%")
	}
}

// LikeInsens filter
func LikeInsens[T FilteredString](field string, value T) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("? ILIKE ?", bun.Ident(field), value+"%")
	}
}

// Page filter
func Page(page Pager) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		if page != nil {
			// limit offset
			if after := page.GetAfter(); after > 0 {
				query = query.Limit(int(after))
			}

			// cursor
			if first := page.GetFirst(); first != "" {
				query = query.Where("id > ?", first)
			}
		}
		return query
	}
}

// Order
func Order(fields ...string) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Order(fields...)
	}
}

// Limit filter
func Limit(limit int) FilterOperator {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Limit(limit)
	}
}

func FilterMerge(query *bun.SelectQuery, filters ...FilterOperator) *bun.SelectQuery {
	for _, filter := range filters {
		if filter != nil {
			query = filter(query)
		}
	}
	return query
}
