package types

import (
	ae "appengine"
)

// Kind represents an entity category in the datastore.
type Kind struct {
	Context ae.Context
	Name    string
}

// NewKind creates a new kind.
func NewKind(ctx ae.Context, name string) *Kind {
	return &Kind{ctx, name}
}
