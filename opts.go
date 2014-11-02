package hrd

type operationOpts struct {
	// completeKeys is whether an entity's key must be set before writing.
	completeKeys bool

	// txCrossGroup is whether the transaction can cross multiple entity groups.
	txCrossGroup bool

	// useLocalCache is whether the in-memory cache is used.
	// useLocalCache bool

	// useGlobalCache is whether memcache is used.
	useGlobalCache bool
}

// Opt is an option to customize the default behaviour of datastore operations.
type Opt int

const (
	// NoCache prevents reading/writing entities from/to the in-memory cache.
	NoCache Opt = iota

	// CompleteKeys prevents entity's key must be set before writing.
	CompleteKeys

	// NoLocalCache prevents reading/writing entities from/to the in-memory cache.
	// NoLocalCache

	// NoGlobalCache prevents reading/writing entities from/to memcache.
	NoGlobalCache
)

func defaultOperationOpts() *operationOpts {
	return &operationOpts{
		useGlobalCache: true,
		//useLocalCache:  true,
	}
}

// clone returns a deep copy.
func (opts *operationOpts) clone() *operationOpts {
	copy := *opts
	return &copy
}

// Flags applies a sequence of Flag.
func (opts *operationOpts) Apply(flags ...Opt) (ret *operationOpts) {
	if len(flags) == 0 {
		return opts
	}

	ret = opts.clone()
	for _, f := range flags {
		switch f {
		case CompleteKeys:
			ret = ret.CompleteKeys(true)
		case NoCache:
			ret = ret.NoCache()
		//case NoLocalCache:
		//	ret = ret.NoLocalCache()
		case NoGlobalCache:
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
	ret.txCrossGroup = true
	return ret
}

// NoCache prevents reading/writing entities from/to the in-memory cache.
func (opts *operationOpts) NoCache() *operationOpts {
	return opts.NoGlobalCache()
}

// NoLocalCache prevents reading/writing entities from/to the in-memory cache.
//func (opts *operationOpts) NoLocalCache() *operationOpts {
//	return opts.NoLocalCacheWrite().NoLocalCacheRead()
//}

// NoGlobalCache prevents reading/writing entities from/to memcache.
func (opts *operationOpts) NoGlobalCache() (ret *operationOpts) {
	ret = opts.clone()
	ret.useGlobalCache = false
	return ret
}

// GlobalCache enables reading/writing entities from/to memcache.
func (opts *operationOpts) GlobalCache() (ret *operationOpts) {
	ret = opts.clone()
	ret.useGlobalCache = true
	return ret
}
