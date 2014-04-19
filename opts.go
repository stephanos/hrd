package hrd

import "time"

type operationOpts struct {
	// completeKeys is whether an entity's key must be set before writing.
	completeKeys bool

	// tx_cross_group is whether the transaction can cross multiple entity groups.
	tx_cross_group bool

	// readLocalCache is whether the in-memory cache is read from.
	readLocalCache bool

	// writeLocalCache is whether the in-memory cache is written to.
	writeLocalCache bool

	// readGlobalCache is whether memcache is read from.
	readGlobalCache bool

	// writeGlobalCache defines an entity's expiration date
	// in memcache. A negative value skips the write.
	writeGlobalCache time.Duration
}

type Opt int

const (
	NO_CACHE Opt = iota
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

// clone returns a deep copy.
func (opts *operationOpts) clone() *operationOpts {
	copy := *opts
	return &copy
}

// Flags applies a sequence of Flag.
func (opts *operationOpts) Apply(flags ...Opt) (ret *operationOpts) {
	ret = opts.clone()
	for _, f := range flags {
		switch f {
		case NO_CACHE:
			ret = ret.NoCache()
		case NO_LOCAL_CACHE:
			ret = ret.NoLocalCache()
		case NO_GLOBAL_CACHE:
			ret = ret.NoGlobalCache()
		}
	}
	return
}

// CompleteKeys defines whether an entity requires a complete key.
// If no parameter is passed, true is assumed.
func (opts *operationOpts) CompleteKeys(complete ...bool) (ret *operationOpts) {
	ret = opts.clone()
	if len(complete) == 1 {
		ret.completeKeys = complete[0]
	} else {
		ret.completeKeys = true
	}
	return ret
}

// XG is whether the transaction can cross multiple entity groups.
func (opts *operationOpts) XG() (ret *operationOpts) {
	ret = opts.clone()
	ret.tx_cross_group = true
	return ret
}

// NoCache prevents reading/writing entities from/to the in-memory cache.
func (opts *operationOpts) NoCache() *operationOpts {
	return opts.NoLocalCache().NoGlobalCache()
}

// NoLocalCache prevents reading/writing entities from/to the in-memory cache.
func (opts *operationOpts) NoLocalCache() *operationOpts {
	return opts.NoLocalCacheWrite().NoLocalCacheRead()
}

// NoGlobalCache prevents reading/writing entities from/to memcache.
func (opts *operationOpts) NoGlobalCache() *operationOpts {
	return opts.NoGlobalCacheWrite().NoGlobalCacheRead()
}

// CacheExpire sets the expiration time for entities written to memcache.
func (opts *operationOpts) CacheExpire(exp time.Duration) (ret *operationOpts) {
	ret = opts.clone()
	ret.writeGlobalCache = exp
	return ret
}

// NoCacheRead prevents reading entities from the in-memory cache or memcache.
func (opts *operationOpts) NoCacheRead() *operationOpts {
	return opts.NoGlobalCacheRead().NoLocalCacheRead()
}

// NoLocalCacheRead prevents reading entities from the in-memory cache.
func (opts *operationOpts) NoLocalCacheRead() (ret *operationOpts) {
	ret = opts.clone()
	ret.readLocalCache = false
	return
}

// NoGlobalCacheRead prevents reading entities from memcache.
func (opts *operationOpts) NoGlobalCacheRead() (ret *operationOpts) {
	ret = opts.clone()
	ret.readGlobalCache = false
	return
}

// NoCacheWrite prevents writing entities to the in-memory cache or memcache.
func (opts *operationOpts) NoCacheWrite() *operationOpts {
	return opts.NoGlobalCacheWrite().NoLocalCacheWrite()
}

// NoLocalCacheWrite prevents writing entities to the in-memory cache.
func (opts *operationOpts) NoLocalCacheWrite() (ret *operationOpts) {
	ret = opts.clone()
	ret.writeLocalCache = false
	return
}

// NoGlobalCacheWrite prevents writing entities to memcache.
func (opts *operationOpts) NoGlobalCacheWrite() *operationOpts {
	return opts.CacheExpire(-1)
}
