package hrd

import "time"

type operationOpts struct {
	completeKeys bool

	tx_cross_group bool

	readLocalCache  bool
	writeLocalCache bool

	readGlobalCache  bool
	writeGlobalCache time.Duration // -1 == skip
}

type Flag int

const (
	NO_CACHE Flag = iota
	NO_LOCAL_CACHE
	NO_GLOBAL_CACHE
)

func defaultOperationOpts() *operationOpts {
	return &operationOpts{
		readLocalCache:   true,
		writeLocalCache:  true,
		readGlobalCache:  true,
		writeGlobalCache: 0,
	}
}

func (self *operationOpts) clone() *operationOpts {
	opts := *self
	return &opts
}

func (self *operationOpts) Flags(flags ...Flag) *operationOpts {
	copy := self.clone()
	for _, f := range flags {
		switch f {
		case NO_CACHE:
			copy = copy.NoCache()
		case NO_LOCAL_CACHE:
			copy = copy.NoLocalCache()
		case NO_GLOBAL_CACHE:
			copy = copy.NoGlobalCache()
		}
	}
	return copy
}

func (self *operationOpts) CompleteKeys(complete ...bool) *operationOpts {
	opts := self.clone()
	if len(complete) == 1 {
		opts.completeKeys = complete[0]
	} else {
		opts.completeKeys = true
	}
	return opts
}

func (self *operationOpts) XG() *operationOpts {
	copy := self.clone()
	copy.tx_cross_group = true
	return copy
}

func (self *operationOpts) NoCache() *operationOpts {
	return self.NoLocalCache().NoGlobalCache()
}

func (self *operationOpts) NoLocalCache() *operationOpts {
	return self.NoLocalCacheWrite().NoLocalCacheRead()
}

func (self *operationOpts) NoGlobalCache() *operationOpts {
	return self.NoGlobalCacheWrite().NoGlobalCacheRead()
}

func (self *operationOpts) CacheExpire(exp time.Duration) *operationOpts {
	copy := self.clone()
	copy.writeGlobalCache = exp
	return copy
}

func (self *operationOpts) NoCacheRead() *operationOpts {
	return self.NoGlobalCacheRead().NoLocalCacheRead()
}

func (self *operationOpts) NoLocalCacheRead() *operationOpts {
	copy := self.clone()
	copy.readLocalCache = false
	return copy
}

func (self *operationOpts) NoGlobalCacheRead() *operationOpts {
	copy := self.clone()
	copy.readGlobalCache = false
	return copy
}

func (self *operationOpts) NoCacheWrite() *operationOpts {
	return self.NoGlobalCacheWrite().NoLocalCacheWrite()
}

func (self *operationOpts) NoLocalCacheWrite() *operationOpts {
	copy := self.clone()
	copy.writeLocalCache = false
	return copy
}

func (self *operationOpts) NoGlobalCacheWrite() *operationOpts {
	return self.CacheExpire(-1)
}
