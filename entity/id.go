package entity

// TextIdentifier can identify a datastore entity via string ID.
type TextIdentifier interface {

	// ID returns a string identifier.
	ID() string

	// SetID applies a string identifier.
	SetID(string)
}

// NumIdentifier can identify a datastore entity via numeric ID.
type NumIdentifier interface {

	// ID returns a numeric identifier.
	ID() int64

	// SetID applies a numeric identifier.
	SetID(int64)
}

// TextParent can identify a datastore entity's parent via string ID.
type TextParent interface {

	// Parent returns the parent's string identifier.
	Parent() string

	// SetParent applies the parent's string identifier.
	SetParent(string)

	// ParentKind returns the parent's collection type.
	ParentKind() string
}

// NumParent can identify a datastore entity's parent via numeric ID.
type NumParent interface {

	// Parent returns the parent's numeric identifier.
	Parent() int64

	// SetParent applies the parent's numeric identifier.
	SetParent(int64)

	// ParentKind returns the parent's collection type.
	ParentKind() string
}
