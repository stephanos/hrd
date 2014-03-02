package hrd

import "time"

type Collection struct {
	opts  *operationOpts
	store *Store
	name  string
}

func (self *Collection) NewNumKey(id int64, parent ...*Key) *Key {
	return self.store.NewNumKey(self.name, id, parent...)
}

func (self *Collection) NewNumKeys(ids ...int64) []*Key {
	return self.store.NewNumKeys(self.name, ids...)
}

func (self *Collection) NewTextKey(id string, parent ...*Key) *Key {
	return self.store.NewTextKey(self.name, id, parent...)
}

func (self *Collection) NewTextKeys(ids ...string) []*Key {
	return self.store.NewTextKeys(self.name, ids...)
}

func (self *Collection) Store() *Store {
	return self.store
}

func (self *Collection) Save() *Saver {
	return newSaver(self)
}

func (self *Collection) Load() *Loader {
	return newLoader(self)
}

func (self *Collection) Delete() *Deleter {
	return &Deleter{self}
}

func (self *Collection) Query() *Query {
	return newQuery(self)
}

// Deletes all entities of the collection.
func (self *Collection) DESTROY() error {
	i := 0
	var start string
	for {
		keys, cursor, err := self.Query().Limit(1000).Start(start).GetKeys()
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			self.store.ctx.Infof("destroyed collection %q (%d items)", self.name, i)
			return nil
		}

		err = self.Delete().Keys(keys)
		if err != nil {
			return err
		}

		start = cursor
		i += len(keys)
	}
}

func (self *Collection) getKey(src interface{}) (*Key, error) {
	return self.store.getKey(self.name, src)
}

// ==== CACHE

func (self *Collection) NoCache() *Collection {
	return self.NoLocalCache().NoGlobalCache()
}

func (self *Collection) NoLocalCache() *Collection {
	return self.NoLocalCacheWrite().NoLocalCacheRead()
}

func (self *Collection) NoGlobalCache() *Collection {
	return self.NoGlobalCacheWrite().NoGlobalCacheRead()
}

func (self *Collection) CacheExpire(exp time.Duration) *Collection {
	self.opts = self.opts.CacheExpire(exp)
	return self
}

func (self *Collection) NoCacheRead() *Collection {
	return self.NoGlobalCacheRead().NoLocalCacheRead()
}

func (self *Collection) NoLocalCacheRead() *Collection {
	self.opts = self.opts.NoLocalCacheRead()
	return self
}

func (self *Collection) NoGlobalCacheRead() *Collection {
	self.opts = self.opts.NoGlobalCacheRead()
	return self
}

func (self *Collection) NoCacheWrite() *Collection {
	return self.NoGlobalCacheWrite().NoLocalCacheWrite()
}

func (self *Collection) NoLocalCacheWrite() *Collection {
	self.opts = self.opts.NoLocalCacheWrite()
	return self
}

func (self *Collection) NoGlobalCacheWrite() *Collection {
	self.opts = self.opts.NoGlobalCacheWrite()
	return self
}
