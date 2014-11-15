package internal

import "appengine"

// Kind represents a category for entities.
type Kind interface {
	Context() appengine.Context

	Name() string
}
