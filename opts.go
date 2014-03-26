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

func (opts *operationOpts) clone() *operationOpts {
	copy := *opts
	return &copy
}

func (opts *operationOpts) Flags(flags ...Flag) (ret *operationOpts) {
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

func (opts *operationOpts) CompleteKeys(complete ...bool) (ret *operationOpts) {
	ret = opts.clone()
	if len(complete) == 1 {
		ret.completeKeys = complete[0]
	} else {
		ret.completeKeys = true
	}
	return ret
}

func (opts *operationOpts) XG() (ret *operationOpts) {
	ret = opts.clone()
	ret.tx_cross_group = true
	return ret
}

func (opts *operationOpts) NoCache() *operationOpts {
	return opts.NoLocalCache().NoGlobalCache()
}

func (opts *operationOpts) NoLocalCache() *operationOpts {
	return opts.NoLocalCacheWrite().NoLocalCacheRead()
}

func (opts *operationOpts) NoGlobalCache() *operationOpts {
	return opts.NoGlobalCacheWrite().NoGlobalCacheRead()
}

func (opts *operationOpts) CacheExpire(exp time.Duration) (ret *operationOpts) {
	ret = opts.clone()
	ret.writeGlobalCache = exp
	return ret
}

func (opts *operationOpts) NoCacheRead() *operationOpts {
	return opts.NoGlobalCacheRead().NoLocalCacheRead()
}

func (opts *operationOpts) NoLocalCacheRead() (ret *operationOpts) {
	ret = opts.clone()
	ret.readLocalCache = false
	return
}

func (opts *operationOpts) NoGlobalCacheRead() (ret *operationOpts) {
	ret = opts.clone()
	ret.readGlobalCache = false
	return
}

func (opts *operationOpts) NoCacheWrite() *operationOpts {
	return opts.NoGlobalCacheWrite().NoLocalCacheWrite()
}

func (opts *operationOpts) NoLocalCacheWrite() (ret *operationOpts) {
	ret = opts.clone()
	ret.writeLocalCache = false
	return
}

func (opts *operationOpts) NoGlobalCacheWrite() *operationOpts {
	return opts.CacheExpire(-1)
}
