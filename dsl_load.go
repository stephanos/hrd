package hrd

import "time"

type Loader struct {
	coll *Collection
	opts *operationOpts
	keys []*Key
}

type SingleLoader struct {
	loader *Loader
}

type MultiLoader struct {
	loader *Loader
}

func newLoader(coll *Collection) *Loader {
	return &Loader{coll: coll, opts: coll.opts.clone()}
}

// Load a single entity by key from the datastore.
func (l *Loader) Key(key *Key) *SingleLoader {
	l.keys = []*Key{key}
	return &SingleLoader{l}
}

// Load multiple entities by key from the datastore.
func (l *Loader) Keys(keys ...*Key) *MultiLoader {
	l.keys = keys
	return &MultiLoader{l}
}

// Load a single entity by id from the datastore.
func (l *Loader) ID(id int64, parent ...*Key) *SingleLoader {
	l.keys = []*Key{l.coll.NewNumKey(id, parent...)}
	return &SingleLoader{l}
}

// Load a single key by text id from the datastore.
func (l *Loader) TextID(id string, parent ...*Key) *SingleLoader {
	l.keys = []*Key{l.coll.NewTextKey(id, parent...)}
	return &SingleLoader{l}
}

// Load multiple keys by id from the datastore.
func (l *Loader) IDs(ids ...int64) *MultiLoader {
	l.keys = l.coll.NewNumKeys(ids...)
	return &MultiLoader{l}
}

// Load multiple keys by text id from the datastore.
func (l *Loader) TextIDs(ids ...string) *MultiLoader {
	l.keys = l.coll.NewTextKeys(ids...)
	return &MultiLoader{l}
}

func (l *Loader) Flags(flags ...Flag) *Loader {
	l.opts = l.opts.Flags(flags...)
	return l
}

// ==== CACHE

func (l *Loader) NoCache() *Loader {
	return l.NoLocalCache().NoGlobalCache()
}

func (l *Loader) NoLocalCache() *Loader {
	return l.NoLocalCacheWrite().NoLocalCacheRead()
}

func (l *Loader) NoGlobalCache() *Loader {
	return l.NoGlobalCacheWrite().NoGlobalCacheRead()
}

func (l *Loader) CacheExpire(exp time.Duration) *Loader {
	l.opts = l.opts.CacheExpire(exp)
	return l
}

func (l *Loader) NoCacheRead() *Loader {
	return l.NoGlobalCacheRead().NoLocalCacheRead()
}

func (l *Loader) NoLocalCacheRead() *Loader {
	l.opts = l.opts.NoLocalCacheRead()
	return l
}

func (l *Loader) NoGlobalCacheRead() *Loader {
	l.opts = l.opts.NoGlobalCacheRead()
	return l
}

func (l *Loader) NoCacheWrite() *Loader {
	return l.NoGlobalCacheWrite().NoLocalCacheWrite()
}

func (l *Loader) NoLocalCacheWrite() *Loader {
	l.opts = l.opts.NoLocalCacheWrite()
	return l
}

func (l *Loader) NoGlobalCacheWrite() *Loader {
	l.opts = l.opts.NoGlobalCacheWrite()
	return l
}

// ==== EXECUTE

// TODO: func (l *Loader) GetEntity(dst interface{}) ([]*Key, error)
// TODO: func (l *Loader) GetEntities(dsts interface{}) ([]*Key, error)

func (l *MultiLoader) GetAll(dsts interface{}) ([]*Key, error) {
	return l.loader.get(dsts, true)
}

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
