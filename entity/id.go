package entity

// TextIdentifier identifies a datastore entity via string ID.
type TextIdentifier interface {

	// ID returns the string identifier.
	ID() string

	// SetID sets the string identifier.
	SetID(string)
}

// TextID represents an entity's string identifier.
type TextID struct {
	Identity string `datastore:"id,inline"`
}

// ID returns the string identifier.
func (mdl *TextID) ID() string {
	return mdl.Identity
}

// SetID sets the string identifier.
func (mdl *TextID) SetID(id string) {
	mdl.Identity = id
}

// NumIdentifier identifies a datastore entity via numeric ID.
type NumIdentifier interface {

	// ID returns the numeric identifier.
	ID() int64

	// SetID sets the numeric identifier.
	SetID(int64)
}

// NumID represents an entity's numeric identifier.
type NumID struct {
	Identity int64 `datastore:"id,inline"`
}

// ID returns the numeric identifier.
func (mdl *NumID) ID() int64 {
	return mdl.Identity
}

// SetID sets the numeric identifier.
func (mdl *NumID) SetID(id int64) {
	mdl.Identity = id
}
