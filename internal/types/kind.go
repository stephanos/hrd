package types

import (
	ae "appengine"
)

// Kind represents a kind of entity in the datastore.
type Kind struct {
	Context ae.Context
	Name    string
}

// NewKind creates a new kind.
func NewKind(ctx ae.Context, name string) *Kind {
	return &Kind{ctx, name}
}
