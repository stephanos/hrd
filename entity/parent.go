package entity

// ParentTextIdentifier identifies an entity's parent via string ID.
type ParentTextIdentifier interface {

	// Parent returns the parent identifier.
	Parent() (kind string, id string)

	// SetParent sets the parent identifier.
	SetParent(kind string, id string)
}

// ParentTextID represents an entity's parent string identifier.
type ParentTextID struct {
	parentKind string
	parentID   string
}

// Parent returns the parent identifier.
func (mdl *ParentTextID) Parent() (kind string, id string) {
	return mdl.parentKind, mdl.parentID
}

// SetParent sets the parent identifier.
func (mdl *ParentTextID) SetParent(kind string, id string) {
	mdl.parentKind = kind
	mdl.parentID = id
}

// ParentNumIdentifier identifies an entity's parent via numeric ID.
type ParentNumIdentifier interface {

	// Parent returns the parent identifier.
	Parent() (kind string, id int64)

	// SetParent sets the parent identifier.
	SetParent(kind string, id int64)
}

// ParentNumID represents an entity's parent numeric identifier.
type ParentNumID struct {
	parentKind string
	parentID   int64
}

// Parent returns the parent identifier.
func (mdl *ParentNumID) Parent() (kind string, id int64) {
	return mdl.parentKind, mdl.parentID
}

// SetParent sets the parent identifier.
func (mdl *ParentNumID) SetParent(kind string, id int64) {
	mdl.parentKind = kind
	mdl.parentID = id
}
