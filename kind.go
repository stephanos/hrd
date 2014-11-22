package hrd

import ae "appengine"

// Kind represents a entity category in the datastore.
type Kind struct {
	store *Store
	name  string
	opts  *opts
}

func newKind(store *Store, name string) *Kind {
	return &Kind{
		store: store,
		name:  name,
		opts:  store.opts.clone(),
	}
}

// Name returns the name of the kind.
func (k *Kind) Name() string {
	return k.name
}

// Save returns a Saver action object.
// It allows to save entities to the datastore.
func (k *Kind) Save(ctx ae.Context) *Saver {
	return newSaver(ctx, k)
}

// Load returns a Loader action object.
// It allows to load entities from the datastore.
func (k *Kind) Load(ctx ae.Context) *Loader {
	return newLoader(ctx, k)
}

// Delete returns a Deleter action object.
// It allows to delete entities from the datastore.
func (k *Kind) Delete(ctx ae.Context) *Deleter {
	return newDeleter(ctx, k)
}

// Query returns a Query object
// It allows to query entities from the datastore.
func (k *Kind) Query(ctx ae.Context) *Query {
	return newQuery(ctx, k)
}

// NewNumKey returns a key for the passed kind and numeric ID.
// It can also receive an optional parent key.
func (k *Kind) NewNumKey(id int64, parent ...*Key) *Key {
	var parentKey *Key
	if len(parent) > 0 {
		parentKey = parent[0]
	}
	return newNumKey(k, id, parentKey)
}

// NewNumKeys returns a sequence of key for the passed kind and
// sequence of numeric ID.
func (k *Kind) NewNumKeys(ids ...int64) []*Key {
	keys := make([]*Key, len(ids))
	for i, id := range ids {
		keys[i] = newNumKey(k, id, nil)
	}
	return keys
}

// NewTextKey returns a key for the passed kind and string ID.
// It can also receive an optional parent key.
func (k *Kind) NewTextKey(id string, parent ...*Key) *Key {
	var parentKey *Key
	if len(parent) > 0 {
		parentKey = parent[0]
	}
	return newTextKey(k, id, parentKey)
}

// NewTextKeys returns a sequence of keys for the passed kind and
// sequence of string ID.
func (k *Kind) NewTextKeys(ids ...string) []*Key {
	keys := make([]*Key, len(ids))
	for i, id := range ids {
		keys[i] = newTextKey(k, id, nil)
	}
	return keys
}
