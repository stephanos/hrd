package hrd

import "time"

type Transactor struct {
	store *Store
	opts  *operationOpts
}

func newTransactor(store *Store) *Transactor {
	return &Transactor{store, store.opts.clone()}
}

func (tx *Transactor) Flags(flags ...Flag) *Transactor {
	tx.opts = tx.opts.Flags(flags...)
	return tx
}

// The transaction can cross multiple entity groups.
func (tx *Transactor) XG() *Transactor {
	tx.opts = tx.opts.XG()
	return tx
}

// ==== CACHE

func (tx *Transactor) NoCache() *Transactor {
	return tx.NoLocalCache().NoGlobalCache()
}

func (tx *Transactor) NoLocalCache() *Transactor {
	return tx.NoLocalCacheWrite().NoLocalCacheRead()
}

func (tx *Transactor) NoGlobalCache() *Transactor {
	return tx.NoGlobalCacheWrite().NoGlobalCacheRead()
}

func (tx *Transactor) CacheExpire(exp time.Duration) *Transactor {
	tx.opts = tx.opts.CacheExpire(exp)
	return tx
}

func (tx *Transactor) NoCacheRead() *Transactor {
	return tx.NoGlobalCacheRead().NoLocalCacheRead()
}

func (tx *Transactor) NoLocalCacheRead() *Transactor {
	tx.opts = tx.opts.NoLocalCacheRead()
	return tx
}

func (tx *Transactor) NoGlobalCacheRead() *Transactor {
	tx.opts = tx.opts.NoGlobalCacheRead()
	return tx
}

func (tx *Transactor) NoCacheWrite() *Transactor {
	return tx.NoGlobalCacheWrite().NoLocalCacheWrite()
}

func (tx *Transactor) NoLocalCacheWrite() *Transactor {
	tx.opts = tx.opts.NoLocalCacheWrite()
	return tx
}

func (tx *Transactor) NoGlobalCacheWrite() *Transactor {
	tx.opts = tx.opts.NoGlobalCacheWrite()
	return tx
}

// ==== EXECUTE

func (tx *Transactor) Run(f func(*Store) ([]*Key, error)) (keys []*Key, err error) {
	return tx.store.runTX(f, tx.opts)
}
