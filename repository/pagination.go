package repository

import "github.com/kabelsea-sandbox/kun/model"

type CursorDirection = int

const (
	CursorNext CursorDirection = iota
	CursorPrev
)

type Cursor = string

func StartCursor[T model.Model](entries []T) Cursor {
	if len(entries) > 0 {
		return Cursor(entries[0].GetID())
	}
	return ""
}

func EndCursor[T model.Model](entries []T) Cursor {
	if len(entries) > 0 {
		return Cursor(entries[len(entries)-1].GetID())
	}
	return ""
}

type Pager interface {
	GetFirst() string
	GetAfter() int32
	GetSkipTotal() bool
}

// PageInfo type
type PageInfo struct {
	Start       Cursor
	End         Cursor
	HasNext     bool
	HasPrevious bool
	Total       int
}
