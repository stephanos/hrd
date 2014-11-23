package internal

import (
	ae "appengine"
	ds "appengine/datastore"
)

// Transact runs a function in a transaction.
func Transact(ctx ae.Context, crossGroup bool, f func(_ ae.Context) error) error {
	return ds.RunInTransaction(ctx, func(ctx ae.Context) error {
		var dsErr error
		dsErr = f(ctx)
		return dsErr
	}, &ds.TransactionOptions{XG: crossGroup})
}
