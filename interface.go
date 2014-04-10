package hrd

import "time"

// textIdentifier can identify a datastore entity via string ID.
type textIdentifier interface {

	// ID returns a string identifier.
	ID() string

	// SetID applies a string identifier.
	SetID(string)
}

// numIdentifier can identify a datastore entity via numeric ID.
type numIdentifier interface {

	// ID returns a numeric identifier.
	ID() int64

	// SetID applies a numeric identifier.
	SetID(int64)
}

// textParent can identify a datastore entity's parent via string ID.
type textParent interface {

	// Parent returns the parent's string identifier.
	Parent() string

	// SetParent applies the parent's string identifier.
	SetParent(string)

	// ParentKind returns the parent's collection type.
	ParentKind() string
}

// numParent can identify a datastore entity's parent via numeric ID.
type numParent interface {

	// Parent returns the parent's numeric identifier.
	Parent() int64

	// SetParent applies the parent's numeric identifier.
	SetParent(int64)

	// ParentKind returns the parent's collection type.
	ParentKind() string
}

// versioned can specify an entity's version.
// It can be used to ignore old cached versions of an entity.
type versioned interface {

	// Version returns a numeric version number.
	Version() int64
}

// timestampCreator can apply the time of creation to an entity.
type timestampCreator interface {

	// SetCreatedAt applies the time of creation.
	SetCreatedAt(time.Time)
}

// timestampUpdater can apply the time of last update to an entity.
type timestampUpdater interface {

	// SetUpdatedAt applies the time of the last update.
	SetUpdatedAt(time.Time)
}

// beforeSaver is a lifecycle hook running before saving an entity.
type beforeSaver interface {

	// BeforeSave runs before an entity is saved.
	// If it returns an error, the save is aborted!
	BeforeSave() error
}

// afterSaver is a lifecycle hook running after saving an entity.
type afterSaver interface {

	// AfterSave runs after an entity is saved.
	AfterSave() error
}

// beforeLoader is a lifecycle hook running before loading an entity.
type beforeLoader interface {

	// BeforeLoad runs before an entity is loaded.
	BeforeLoad() error
}

// afterLoader is a lifecycle hook running after loading an entity.
type afterLoader interface {

	// AfterLoad runs after an entity is loaded.
	AfterLoad() error
}
