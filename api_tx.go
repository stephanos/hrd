package hrd

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

// NoGlobalCache prevents reading/writing entities from/to memcache.
func (tx *Transactor) NoGlobalCache() *Transactor {
	tx.opts = tx.opts.NoGlobalCache()
	return tx
}

// GlobalCache enables reading/writing entities from/to memcache.
func (tx *Transactor) GlobalCache() *Transactor {
	tx.opts = tx.opts.GlobalCache()
	return tx
}

// ==== EXECUTE

// Run executes a transaction.
func (tx *Transactor) Run(f func(*Store) error) error {
	return tx.store.runTX(f, tx.opts)
}
