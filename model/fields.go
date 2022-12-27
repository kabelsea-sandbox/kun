package model

import (
	"time"
)

// Field interface
type Field interface {
	IsModelField()
}

// ID
type ID = string

// DateTime
type DateTime = time.Time
