package hrd

import "time"

// Transactor can run multiple datastore operations inside a transaction.
type Transactor struct {
	store *Store
	opts  *operationOpts
}

// newTransactor creates a new Transactor for the passed store.
// The store's options are used as default options.
func newTransactor(store *Store) *Transactor {
	return &Transactor{store, store.opts.clone()}
}

// Opts applies the sequence of Opt to the Transactor's options.
func (tx *Transactor) Opts(opts ...Opt) *Transactor {
	tx.opts = tx.opts.Apply(opts...)
	return tx
}

// XG allows the transaction to run across multiple entity groups.
func (tx *Transactor) XG() *Transactor {
	tx.opts = tx.opts.XG()
	return tx
}

// ==== CACHE

// NoCache prevents reading/writing entities from/to
// the in-memory cache or memcache in this transaction.
func (tx *Transactor) NoCache() *Transactor {
	return tx.NoLocalCache().NoGlobalCache()
}

// NoLocalCache prevents reading/writing entities from/to
// the in-memory cache in this transaction.
func (tx *Transactor) NoLocalCache() *Transactor {
	return tx.NoLocalCacheWrite().NoLocalCacheRead()
}

// NoGlobalCache prevents writing entities to
// memcache in this transaction.
func (tx *Transactor) NoGlobalCache() *Transactor {
	return tx.NoGlobalCacheWrite().NoGlobalCacheRead()
}

// CacheExpire sets the expiration time in memcache for entities
// that are cached after the successful transaction.
func (tx *Transactor) CacheExpire(exp time.Duration) *Transactor {
	tx.opts = tx.opts.CacheExpire(exp)
	return tx
}

// NoCacheRead prevents reading entities from
// the in-memory cache or memcache in this transaction.
func (tx *Transactor) NoCacheRead() *Transactor {
	return tx.NoGlobalCacheRead().NoLocalCacheRead()
}

// NoLocalCacheRead prevents reading entities from
// the in-memory cache in this transaction.
func (tx *Transactor) NoLocalCacheRead() *Transactor {
	tx.opts = tx.opts.NoLocalCacheRead()
	return tx
}

// NoGlobalCache prevents writing entities to
// memcache in this transaction.
func (tx *Transactor) NoGlobalCacheRead() *Transactor {
	tx.opts = tx.opts.NoGlobalCacheRead()
	return tx
}

// NoCacheWrite prevents writing entities to
// the in-memory cache or memcache in this transaction.
func (tx *Transactor) NoCacheWrite() *Transactor {
	return tx.NoGlobalCacheWrite().NoLocalCacheWrite()
}

// NoLocalCacheWrite prevents writing entities to
// the in-memory cache in this transaction.
func (tx *Transactor) NoLocalCacheWrite() *Transactor {
	tx.opts = tx.opts.NoLocalCacheWrite()
	return tx
}

// NoGlobalCacheWrite prevents writing entities to
// memcache in this transaction.
func (tx *Transactor) NoGlobalCacheWrite() *Transactor {
	tx.opts = tx.opts.NoGlobalCacheWrite()
	return tx
}

// ==== EXECUTE

// Run executes a transaction.
func (tx *Transactor) Run(f func(*Store) ([]*Key, error)) (keys []*Key, err error) {
	return tx.store.runTX(f, tx.opts)
}
