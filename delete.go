package hrd

import ae "appengine"

// Deleter can delete entities from a kind.
type Deleter struct {
	*actionContext
}

// newDeleter creates a new Deleter for the passed kind.
// The kind's options are used as default options.
func newDeleter(ctx ae.Context, kind *Kind) *Deleter {
	return &Deleter{newActionContext(ctx, kind)}
}

// Key deletes a single entity by key from the datastore.
func (d *Deleter) Key(key *Key) error {
	return d.deleteKeys(key)
}

// Keys deletes multiple entities by key from the datastore.
func (d *Deleter) Keys(keys []*Key) error {
	return d.deleteKeys(keys...)
}

// ID deletes a single entity by id from the datastore.
func (d *Deleter) ID(id int64, parent ...*Key) error {
	return d.deleteKeys(d.kind.NewNumKey(id, parent...))
}

// TextID deletes a single key by text id from the datastore.
func (d *Deleter) TextID(id string, parent ...*Key) error {
	return d.deleteKeys(d.kind.NewTextKey(id, parent...))
}

// IDs deletes multiple keys by id from the datastore.
func (d *Deleter) IDs(ids ...int64) error {
	return d.deleteKeys(d.kind.NewNumKeys(ids...)...)
}

// TextIDs deletes multiple keys by text id from the datastore.
func (d *Deleter) TextIDs(ids ...string) error {
	return d.deleteKeys(d.kind.NewTextKeys(ids...)...)
}

// Entity deletes the provided entity.
func (d *Deleter) Entity(src interface{}) error {
	return dsDelete(d.Kind(), src, false)
}

// Entities deletes the provided entities.
func (d *Deleter) Entities(srcs interface{}) error {
	return dsDelete(d.Kind(), srcs, true)
}

func (d *Deleter) deleteKeys(keys ...*Key) error {
	return dsDeleteKeys(d.Kind(), toInternalKeys(keys)...)
}
