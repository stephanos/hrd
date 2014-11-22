package hrd

import ae "appengine"

// Transactor can run multiple datastore operations inside a transaction.
// By default it does not handle multiple entity groups.
type Transactor struct {
	ctx        ae.Context
	opts       *Opts
	crossGroup bool
}

// TX represents an App Engine Context inside a transaction.
// It should only be used inside a transaction.
type TX interface {
	ae.Context
}

func newTransactor(s *Store, ctx ae.Context) *Transactor {
	return &Transactor{ctx: ctx, opts: s.opts.clone()}
}

// XG defines whether the transaction can cross multiple entity groups.
// If no parameter is passed, true is assumed.
func (tx *Transactor) XG(enable ...bool) *Transactor {
	crossGroup := true
	if len(enable) > 0 {
		crossGroup = enable[0]
	}
	tx.crossGroup = crossGroup
	return tx
}

// Run executes a function in a transaction.
func (tx *Transactor) Run(f func(_ TX) error) error {
	return dsTransact(tx.ctx, func(ctx ae.Context) error {
		return f(ctx)
	}, tx.crossGroup)
}
