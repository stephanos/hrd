package entity

// BeforeSaver is a lifecycle hook running before saving an entity.
type BeforeSaver interface {

	// BeforeSave runs before an entity is saved.
	// If it returns an error, the save is aborted!
	BeforeSave() error
}

// AfterSaver is a lifecycle hook running after saving an entity.
type AfterSaver interface {

	// AfterSave runs after an entity is saved.
	AfterSave() error
}

// BeforeLoader is a lifecycle hook running before loading an entity.
type BeforeLoader interface {

	// BeforeLoad runs before an entity is loaded.
	BeforeLoad() error
}

// AfterLoader is a lifecycle hook running after loading an entity.
type AfterLoader interface {

	// AfterLoad runs after an entity is loaded.
	AfterLoad() error
}
