package internal

import ae "appengine"

// Kind represents a category for entities.
type Kind interface {
	Context() ae.Context

	Name() string
}
