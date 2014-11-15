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

// ==== CONFIG

// Opts applies the sequence of Opt to the Transactor's options.
func (tx *Transactor) Opts(opts ...Opt) *Transactor {
	tx.opts = tx.opts.Apply(opts...)
	return tx
}

// XG defines whether the transaction can cross multiple entity groups.
// If no parameter is passed, true is assumed.
func (tx *Transactor) XG(enable ...bool) *Transactor {
	tx.opts = tx.opts.XG(enable...)
	return tx
}

// ==== EXECUTE

// Run executes a transaction.
func (tx *Transactor) Run(f func(*Store) error) error {
	return tx.store.runTX(f, tx.opts)
}
