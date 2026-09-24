package domain

/*
Nullable needs to specify:
- field not provided
- field provided: null
- field provided: vlaue
*/

type Nullable[T any] struct {
	Value *T
	Set   bool
}
