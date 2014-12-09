package types

// Opts represents options that change the behaviour of datastore operations.
type Opts struct {

	// CompleteKeys is whether an entity's key must be set before writing.
	CompleteKeys bool

	// NoGlobalCache is whether memcache is used.
	NoGlobalCache bool
}

// DefaultOpts returns an object with default options.
func DefaultOpts() *Opts {
	return &Opts{}
}

// Clone returns a deep copy.
func (opts *Opts) Clone() *Opts {
	copy := *opts
	return &copy
}
