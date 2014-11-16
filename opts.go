package hrd

type operationOpts struct {
	// completeKeys is whether an entity's key must be set before writing.
	completeKeys bool

	// txCrossGroup is whether the transaction can cross multiple entity groups.
	txCrossGroup bool

	// useGlobalCache is whether memcache is used.
	useGlobalCache bool
}

// Opt is an option to customize the default behaviour of datastore operations.
type Opt int

const (
	// CompleteKeys prevents entity's key must be set before writing.
	CompleteKeys = iota

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
		case NoGlobalCache:
			ret = ret.GlobalCache(false)
		}
	}
	return
}

// CompleteKeys defines whether an entity requires a complete key.
// If no parameter is passed, true is assumed.
func (opts *operationOpts) CompleteKeys(complete ...bool) (ret *operationOpts) {
	ret = opts.clone()
	ret.completeKeys = allTrueOrEmpty(complete...)
	return ret
}

// XG defines whether the transaction can cross multiple entity groups.
// If no parameter is passed, true is assumed.
func (opts *operationOpts) XG(enable ...bool) (ret *operationOpts) {
	ret = opts.clone()
	ret.txCrossGroup = allTrueOrEmpty(enable...)
	return ret
}

// GlobalCache defines whether entities are read/written from/to memcache.
// If no parameter is passed, true is assumed.
func (opts *operationOpts) GlobalCache(enable ...bool) (ret *operationOpts) {
	ret = opts.clone()
	ret.useGlobalCache = allTrueOrEmpty(enable...)
	return
}

func allTrueOrEmpty(bools ...bool) bool {
	for _, b := range bools {
		if !b {
			return false
		}
	}
	return true
}
