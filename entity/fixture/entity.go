package fixture

import "github.com/101loops/hrd/entity"

// EntityWithNumID is an entity with a numeric identifier.
type EntityWithNumID struct {
	entity.NumID
}

// EntityWithTextID is an entity with a text identifier.
type EntityWithTextID struct {
	entity.TextID
}

// EntityWithParentNumID is an entity with a numeric parent identifier.
type EntityWithParentNumID struct {
	entity.NumID

	parentKind string
	parentID   int64
}

// Parent returns the entity's parent.
func (mdl *EntityWithParentNumID) Parent() (kind string, id int64) {
	return mdl.parentKind, mdl.parentID
}

// SetParent applies the entity's parent.
func (mdl *EntityWithParentNumID) SetParent(kind string, id int64) {
	mdl.parentKind = kind
	mdl.parentID = id
}

// EntityWithParentTextID is an entity with a text parent identifier.
type EntityWithParentTextID struct {
	entity.TextID

	parentKind string
	parentID   string
}

// Parent returns the entity's parent.
func (mdl *EntityWithParentTextID) Parent() (kind string, id string) {
	return mdl.parentKind, mdl.parentID
}

// SetParent applies the entity's parent.
func (mdl *EntityWithParentTextID) SetParent(kind string, id string) {
	mdl.parentKind = kind
	mdl.parentID = id
}
