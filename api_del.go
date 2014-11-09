package hrd

// Deleter can delete entities from a Collection.
type Deleter struct {
	coll *Collection
}

// newDeleter creates a new Deleter for the passed collection.
// The collection's options are used as default options.
func newDeleter(coll *Collection) *Deleter {
	return &Deleter{coll: coll}
}

// ==== EXECUTE

// Key deletes a single entity by key from the datastore.
func (d *Deleter) Key(key *Key) error {
	return d.delete([]*Key{key})
}

// Keys deletes multiple entities by key from the datastore.
func (d *Deleter) Keys(keys []*Key) error {
	return d.delete(keys)
}

// Entity deletes the provided entity.
func (d *Deleter) Entity(src interface{}) error {
	k, err := d.coll.store.getKey(d.coll.name, src)
	if err != nil {
		return err
	}
	return d.delete([]*Key{k})
}

// Entities deletes the provided entities.
func (d *Deleter) Entities(srcs interface{}) error {
	dsKeys, err := d.coll.store.getKeys(d.coll.name, srcs)
	if err != nil {
		return err
	}
	return d.delete(dsKeys)
}

// ID deletes a single entity by id from the datastore.
func (d *Deleter) ID(id int64, parent ...*Key) error {
	return d.delete([]*Key{d.coll.NewNumKey(id, parent...)})
}

// TextID deletes a single key by text id from the datastore.
func (d *Deleter) TextID(id string, parent ...*Key) error {
	return d.delete([]*Key{d.coll.NewTextKey(id, parent...)})
}

// IDs deletes multiple keys by id from the datastore.
func (d *Deleter) IDs(ids ...int64) error {
	return d.delete(d.coll.NewNumKeys(ids...))
}

// TextIDs deletes multiple keys by text id from the datastore.
func (d *Deleter) TextIDs(ids ...string) error {
	return d.delete(d.coll.NewTextKeys(ids...))
}

func (d *Deleter) delete(keys []*Key) error {
	return d.coll.store.deleteMulti(d.coll.name, keys)
}
