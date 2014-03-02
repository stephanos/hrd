package hrd

import "time"

// Datastore ID
type textIdentifier interface {
	ID() string
	SetID(string)
}

type numIdentifier interface {
	ID() int64
	SetID(int64)
}

// Datastore parent ID
type textParent interface {
	Parent() string
	ParentKind() string
	SetParent(string)
}

type numParent interface {
	Parent() int64
	ParentKind() string
	SetParent(int64)
}

type versioned interface {
	Version() int64
}

// Datastore timestamp
type timestampCreator interface {
	SetCreatedAt(time.Time)
}

type timestampUpdater interface {
	SetUpdatedAt(time.Time)
}

// Datastore "save" hooks
type beforeSaver interface {
	BeforeSave() error
}

type afterSaver interface {
	AfterSave() error
}

// Datastore "load" hooks
type beforeLoader interface {
	BeforeLoad() error
}

type afterLoader interface {
	AfterLoad() error
}
