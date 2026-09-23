package core_http_types

import (
	"encoding/json"

	"github.com/apollobing/todo-app/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}

func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil
		
		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set: n.Set,
	}
}

/*
------------------
JSON: {}
Nullable:
	- Value: *nil
	- Set: false


------------------
JSON: {
	"phone_number": "+491234567890"
}
Nullable:
	- Value: *"+491234567890"
	- Set: true


------------------
JSON: {
	"phone_number": null
}
Nullable:
	- Value: *nil
	- Set: true
*/
