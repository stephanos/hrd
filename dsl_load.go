package hrd

import "time"

// Loader can load entities from a Collection.
type Loader struct {
	coll *Collection
	opts *operationOpts
	keys []*Key
}

// SingleLoader is a special Loader that allows to fetch exactly one entity
// from the datastore.
type SingleLoader struct {
	loader *Loader
}

// MultiLoader is a special Loader that allows to fetch multiple entities
// from the datastore.
type MultiLoader struct {
	loader *Loader
}

// newLoader creates a new Loader for the passed collection.
// The collection's options are used as default options.
func newLoader(coll *Collection) *Loader {
	return &Loader{coll: coll, opts: coll.opts.clone()}
}

// Key loads a single entity by key from the datastore.
func (l *Loader) Key(key *Key) *SingleLoader {
	l.keys = []*Key{key}
	return &SingleLoader{l}
}

// Keys load multiple entities by key from the datastore.
func (l *Loader) Keys(keys ...*Key) *MultiLoader {
	l.keys = keys
	return &MultiLoader{l}
}

// ID loads a single entity by id from the datastore.
func (l *Loader) ID(id int64, parent ...*Key) *SingleLoader {
	l.keys = []*Key{l.coll.NewNumKey(id, parent...)}
	return &SingleLoader{l}
}

// TextID loads a single key by text id from the datastore.
func (l *Loader) TextID(id string, parent ...*Key) *SingleLoader {
	l.keys = []*Key{l.coll.NewTextKey(id, parent...)}
	return &SingleLoader{l}
}

// IDs load multiple keys by id from the datastore.
func (l *Loader) IDs(ids ...int64) *MultiLoader {
	l.keys = l.coll.NewNumKeys(ids...)
	return &MultiLoader{l}
}

// TextIDs load multiple keys by text id from the datastore.
func (l *Loader) TextIDs(ids ...string) *MultiLoader {
	l.keys = l.coll.NewTextKeys(ids...)
	return &MultiLoader{l}
}

// Opts applies a sequence of Opt the Loader's options.
func (l *Loader) Opts(opts ...Opt) *Loader {
	l.opts = l.opts.Apply(opts...)
	return l
}

// ==== CACHE

// NoCache prevents reading/writing entities from/to
// the in-memory cache or memcache in this load operation.
func (l *Loader) NoCache() *Loader {
	return l.NoLocalCache().NoGlobalCache()
}

// NoLocalCache prevents reading/writing entities from/to
// the in-memory cache in this load operation.
func (l *Loader) NoLocalCache() *Loader {
	return l.NoLocalCacheWrite().NoLocalCacheRead()
}

// NoGlobalCache prevents reading/writing entities from/to
// memcache in this load operation.
func (l *Loader) NoGlobalCache() *Loader {
	return l.NoGlobalCacheWrite().NoGlobalCacheRead()
}

// CacheExpire sets the expiration time in memcache for entities
// that are cached after loading them from the datastore.
func (l *Loader) CacheExpire(exp time.Duration) *Loader {
	l.opts = l.opts.CacheExpire(exp)
	return l
}

// NoCacheRead prevents reading entities from
// the in-memory cache or memcache in this load operation.
func (l *Loader) NoCacheRead() *Loader {
	return l.NoGlobalCacheRead().NoLocalCacheRead()
}

// NoLocalCacheRead prevents reading entities from
// the in-memory cache in this load operation.
func (l *Loader) NoLocalCacheRead() *Loader {
	l.opts = l.opts.NoLocalCacheRead()
	return l
}

// NoGlobalCacheRead prevents reading entities from
// memcache in this load operation.
func (l *Loader) NoGlobalCacheRead() *Loader {
	l.opts = l.opts.NoGlobalCacheRead()
	return l
}

// NoCacheWrite prevents writing entities to
// the in-memory cache or memcache in this load operation.
func (l *Loader) NoCacheWrite() *Loader {
	return l.NoGlobalCacheWrite().NoLocalCacheWrite()
}

// NoLocalCacheWrite prevents writing entities to
// the in-memory cache in this load operation.
func (l *Loader) NoLocalCacheWrite() *Loader {
	l.opts = l.opts.NoLocalCacheWrite()
	return l
}

// NoGlobalCacheWrite prevents writing entities to
// memcache in this load operation.
func (l *Loader) NoGlobalCacheWrite() *Loader {
	l.opts = l.opts.NoGlobalCacheWrite()
	return l
}

// ==== EXECUTE

// TODO: func (l *Loader) GetEntity(dst interface{}) ([]*Key, error)
// TODO: func (l *Loader) GetEntities(dsts interface{}) ([]*Key, error)

// GetAll loads entities from the datastore into the passed destination.
func (l *MultiLoader) GetAll(dsts interface{}) ([]*Key, error) {
	return l.loader.get(dsts, true)
}

// GetOne loads an entity from the datastore into the passed destination.
func (l *SingleLoader) GetOne(dst interface{}) (*Key, error) {
	var key *Key
	keys, err := l.loader.get(dst, false)
	if len(keys) == 1 {
		if keys[0].Exists() {
			key = keys[0]
		}
	}
	return key, err
}

func (l *Loader) get(dst interface{}, multi bool) ([]*Key, error) {
	docs, err := newWriteableDocs(dst, l.keys, multi)
	if err != nil {
		return nil, err
	}
	return l.coll.store.getMulti(l.coll.name, docs, l.opts)
}
