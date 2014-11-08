package hrd

// Collection represents a datastore kind.
type Collection struct {
	opts  *operationOpts
	store *Store
	name  string
}

// Name returns the name of the collection.
func (coll *Collection) Name() string {
	return coll.name
}

// NewNumKey returns a key for the passed numeric ID.
// It can also receive an optional parent key.
func (coll *Collection) NewNumKey(id int64, parent ...*Key) *Key {
	return coll.store.NewNumKey(coll.name, id, parent...)
}

// NewNumKeys returns a sequence of key for the sequence of numeric ID.
func (coll *Collection) NewNumKeys(ids ...int64) []*Key {
	return coll.store.NewNumKeys(coll.name, ids...)
}

// NewTextKey returns a key for the passed string ID.
// It can also receive an optional parent key.
func (coll *Collection) NewTextKey(id string, parent ...*Key) *Key {
	return coll.store.NewTextKey(coll.name, id, parent...)
}

// NewTextKeys returns a sequence of keys for the passed sequence of string ID.
func (coll *Collection) NewTextKeys(ids ...string) []*Key {
	return coll.store.NewTextKeys(coll.name, ids...)
}

// Save returns a Saver action object.
// It allows to save entities to the datastore.
func (coll *Collection) Save(opts ...Opt) *Saver {
	return newSaver(coll).Opts(opts...)
}

// Load returns a Loader action object.
// It allows to load entities from the datastore.
func (coll *Collection) Load(opts ...Opt) *Loader {
	return newLoader(coll).Opts(opts...)
}

// Delete returns a Deleter action object.
// It allows to delete entities from the datastore.
func (coll *Collection) Delete() *Deleter {
	return newDeleter(coll)
}

// Query returns a Query object
// It allows to query entities from the datastore.
func (coll *Collection) Query(opts ...Opt) *Query {
	return newQuery(coll).Opts(opts...)
}

// DESTROY deletes all entities of the collection.
// Proceed with extreme caution!
func (coll *Collection) DESTROY() ([]*Key, error) {
	var i int
	var start string
	var allKeys []*Key
	for {
		keys, cursor, dsErr := coll.Query().Limit(1000).Start(start).GetKeys()
		if dsErr != nil {
			return allKeys, dsErr
		}
		if len(keys) == 0 {
			coll.store.ctx.Infof("destroyed collection %q (%d items)", coll.name, i)
			return allKeys, nil
		}

		dsErr = coll.Delete().Keys(keys)
		if dsErr != nil {
			return allKeys, dsErr
		}

		start = cursor
		i += len(keys)
		allKeys = append(allKeys, keys...)
	}
}

func (coll *Collection) getKey(src interface{}) (*Key, error) {
	return coll.store.getKey(coll.name, src)
}
