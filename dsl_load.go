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
func (self *Loader) Key(key *Key) *SingleLoader {
	self.keys = []*Key{key}
	return &SingleLoader{self}
}

// Load multiple entities by key from the datastore.
func (self *Loader) Keys(keys ...*Key) *MultiLoader {
	self.keys = keys
	return &MultiLoader{self}
}

// Load a single entity by id from the datastore.
func (self *Loader) ID(id int64, parent ...*Key) *SingleLoader {
	self.keys = []*Key{self.coll.NewNumKey(id, parent...)}
	return &SingleLoader{self}
}

// Load a single key by text id from the datastore.
func (self *Loader) TextID(id string, parent ...*Key) *SingleLoader {
	self.keys = []*Key{self.coll.NewTextKey(id, parent...)}
	return &SingleLoader{self}
}

// Load multiple keys by id from the datastore.
func (self *Loader) IDs(ids ...int64) *MultiLoader {
	self.keys = self.coll.NewNumKeys(ids...)
	return &MultiLoader{self}
}

// Load multiple keys by text id from the datastore.
func (self *Loader) TextIDs(ids ...string) *MultiLoader {
	self.keys = self.coll.NewTextKeys(ids...)
	return &MultiLoader{self}
}

func (self *Loader) Flags(flags ...Flag) *Loader {
	self.opts = self.opts.Flags(flags...)
	return self
}

// ==== CACHE

func (self *Loader) NoCache() *Loader {
	return self.NoLocalCache().NoGlobalCache()
}

func (self *Loader) NoLocalCache() *Loader {
	return self.NoLocalCacheWrite().NoLocalCacheRead()
}

func (self *Loader) NoGlobalCache() *Loader {
	return self.NoGlobalCacheWrite().NoGlobalCacheRead()
}

func (self *Loader) CacheExpire(exp time.Duration) *Loader {
	self.opts = self.opts.CacheExpire(exp)
	return self
}

func (self *Loader) NoCacheRead() *Loader {
	return self.NoGlobalCacheRead().NoLocalCacheRead()
}

func (self *Loader) NoLocalCacheRead() *Loader {
	self.opts = self.opts.NoLocalCacheRead()
	return self
}

func (self *Loader) NoGlobalCacheRead() *Loader {
	self.opts = self.opts.NoGlobalCacheRead()
	return self
}

func (self *Loader) NoCacheWrite() *Loader {
	return self.NoGlobalCacheWrite().NoLocalCacheWrite()
}

func (self *Loader) NoLocalCacheWrite() *Loader {
	self.opts = self.opts.NoLocalCacheWrite()
	return self
}

func (self *Loader) NoGlobalCacheWrite() *Loader {
	self.opts = self.opts.NoGlobalCacheWrite()
	return self
}

// ==== EXECUTE

// TODO: func (self *Loader) GetEntity(dst interface{}) ([]*Key, error)
// TODO: func (self *Loader) GetEntities(dsts interface{}) ([]*Key, error)

func (self *MultiLoader) GetAll(dsts interface{}) ([]*Key, error) {
	return self.loader.get(dsts, true)
}

func (self *SingleLoader) GetOne(dst interface{}) (*Key, error) {
	var key *Key
	keys, err := self.loader.get(dst, false)
	if len(keys) == 1 {
		if keys[0].Exists() {
			key = keys[0]
		}
	}
	return key, err
}

func (self *Loader) get(dst interface{}, multi bool) ([]*Key, error) {
	docs, err := newWriteableDocs(dst, self.keys, multi)
	if err != nil {
		return nil, err
	}
	return self.coll.store.getMulti(self.coll.name, docs, self.opts)
}
