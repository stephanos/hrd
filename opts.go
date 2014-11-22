package hrd

type Opts struct {
	// completeKeys is whether an entity's key must be set before writing.
	completeKeys bool

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

func defaultOpts() *Opts {
	return &Opts{
		useGlobalCache: true,
	}
}

// Clone returns a deep copy.
func (opts *Opts) clone() *Opts {
	copy := *opts
	return &copy
}

// Flags applies a sequence of Flag.
func (opts *Opts) Apply(flags ...Opt) (ret *Opts) {
	if len(flags) == 0 {
		return opts
	}

	ret = opts.clone()
	for _, f := range flags {
		switch f {
		case CompleteKeys:
			ret.completeKeys = true
		case NoGlobalCache:
			ret.useGlobalCache = false
		}
	}
	return
}
