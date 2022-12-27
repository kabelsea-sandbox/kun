package repository

// Operator interface
type Operator interface{}

type OperatorContext[T Operator] struct {
	operatorContext T
}
