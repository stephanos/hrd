package hrd

import "time"

type Collection struct {
	opts  *operationOpts
	store *Store
	name  string
}

func (coll *Collection) NewNumKey(id int64, parent ...*Key) *Key {
	return coll.store.NewNumKey(coll.name, id, parent...)
}

func (coll *Collection) NewNumKeys(ids ...int64) []*Key {
	return coll.store.NewNumKeys(coll.name, ids...)
}

func (coll *Collection) NewTextKey(id string, parent ...*Key) *Key {
	return coll.store.NewTextKey(coll.name, id, parent...)
}

func (coll *Collection) NewTextKeys(ids ...string) []*Key {
	return coll.store.NewTextKeys(coll.name, ids...)
}

func (coll *Collection) Store() *Store {
	return coll.store
}

func (coll *Collection) Save() *Saver {
	return newSaver(coll)
}

func (coll *Collection) Load() *Loader {
	return newLoader(coll)
}

func (coll *Collection) Delete() *Deleter {
	return &Deleter{coll}
}

func (coll *Collection) Query() *Query {
	return newQuery(coll)
}

// Deletes all entities of the collection.
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

func (coll *Collection) NoCache() *Collection {
	return coll.NoLocalCache().NoGlobalCache()
}

func (coll *Collection) NoLocalCache() *Collection {
	return coll.NoLocalCacheWrite().NoLocalCacheRead()
}

func (coll *Collection) NoGlobalCache() *Collection {
	return coll.NoGlobalCacheWrite().NoGlobalCacheRead()
}

func (coll *Collection) CacheExpire(exp time.Duration) *Collection {
	coll.opts = coll.opts.CacheExpire(exp)
	return coll
}

func (coll *Collection) NoCacheRead() *Collection {
	return coll.NoGlobalCacheRead().NoLocalCacheRead()
}

func (coll *Collection) NoLocalCacheRead() *Collection {
	coll.opts = coll.opts.NoLocalCacheRead()
	return coll
}

func (coll *Collection) NoGlobalCacheRead() *Collection {
	coll.opts = coll.opts.NoGlobalCacheRead()
	return coll
}

func (coll *Collection) NoCacheWrite() *Collection {
	return coll.NoGlobalCacheWrite().NoLocalCacheWrite()
}

func (coll *Collection) NoLocalCacheWrite() *Collection {
	coll.opts = coll.opts.NoLocalCacheWrite()
	return coll
}

func (coll *Collection) NoGlobalCacheWrite() *Collection {
	coll.opts = coll.opts.NoGlobalCacheWrite()
	return coll
}
