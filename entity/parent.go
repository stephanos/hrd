package entity

// ParentTextIdentifier identifies an entity's parent via string ID.
type ParentTextIdentifier interface {

	// Parent returns the parent identifier.
	Parent() (kind string, id string)

	// SetParent sets the parent identifier.
	SetParent(kind string, id string)
}

// ParentNumIdentifier identifies an entity's parent via numeric ID.
type ParentNumIdentifier interface {

	// Parent returns the parent identifier.
	Parent() (kind string, id int64)

	// SetParent sets the parent identifier.
	SetParent(kind string, id int64)
}
