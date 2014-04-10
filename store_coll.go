package hrd

import "time"

type Collection struct {
	opts  *operationOpts
	store *Store
	name  string
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
func (coll *Collection) Save() *Saver {
	return newSaver(coll)
}

// Load returns a Loader action object.
// It allows to load entities from the datastore.
func (coll *Collection) Load() *Loader {
	return newLoader(coll)
}

// Delete returns a Deleter action object.
// It allows to delete entities from the datastore.
func (coll *Collection) Delete() *Deleter {
	return newDeleter(coll)
}

// Query returns a Query object
// It allows to query entities from the datastore.
func (coll *Collection) Query() *Query {
	return newQuery(coll)
}

// DESTROY deletes all entities of the collection.
// Proceed with extreme caution!
func (coll *Collection) DESTROY() error {
	i := 0
	var start string
	for {
		keys, cursor, err := coll.Query().Limit(1000).Start(start).GetKeys()
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			coll.store.ctx.Infof("destroyed collection %q (%d items)", coll.name, i)
			return nil
		}

		err = coll.Delete().Keys(keys)
		if err != nil {
			return err
		}

		start = cursor
		i += len(keys)
	}
}

func (coll *Collection) getKey(src interface{}) (*Key, error) {
	return coll.store.getKey(coll.name, src)
}

// ==== CACHE

// NoCache prevents entities of this collection to be cached in-memory or in memcache.
func (coll *Collection) NoCache() *Collection {
	return coll.NoLocalCache().NoGlobalCache()
}

// NoCache prevents entities of this collection to be cached in-memory.
func (coll *Collection) NoLocalCache() *Collection {
	return coll.NoLocalCacheWrite().NoLocalCacheRead()
}

// NoCache prevents entities of this collection to be cached in memcache.
func (coll *Collection) NoGlobalCache() *Collection {
	return coll.NoGlobalCacheWrite().NoGlobalCacheRead()
}

// CacheExpire sets the expiration time for entities written to memcache.
func (coll *Collection) CacheExpire(exp time.Duration) *Collection {
	coll.opts = coll.opts.CacheExpire(exp)
	return coll
}

// NoCacheRead prevents reading entities from the in-memory cache or memcache.
func (coll *Collection) NoCacheRead() *Collection {
	return coll.NoGlobalCacheRead().NoLocalCacheRead()
}

// NoLocalCacheRead prevents reading entities from the in-memory cache.
func (coll *Collection) NoLocalCacheRead() *Collection {
	coll.opts = coll.opts.NoLocalCacheRead()
	return coll
}

// NoGlobalCacheRead prevents reading entities from memcache.
func (coll *Collection) NoGlobalCacheRead() *Collection {
	coll.opts = coll.opts.NoGlobalCacheRead()
	return coll
}

// NoCacheWrite prevents writing entities to the in-memory cache or memcache.
func (coll *Collection) NoCacheWrite() *Collection {
	return coll.NoGlobalCacheWrite().NoLocalCacheWrite()
}

// NoLocalCacheWrite prevents writing entities to the in-memory cache.
func (coll *Collection) NoLocalCacheWrite() *Collection {
	coll.opts = coll.opts.NoLocalCacheWrite()
	return coll
}

// NoGlobalCacheWrite prevents writing entities to memcache.
func (coll *Collection) NoGlobalCacheWrite() *Collection {
	coll.opts = coll.opts.NoGlobalCacheWrite()
	return coll
}
