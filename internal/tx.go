package internal

import (
	ae "appengine"
	ds "appengine/datastore"
)

// DSTransact runs a function in a transaction.
func DSTransact(ctx ae.Context, f func(_ ae.Context) error, crossGroup bool) error {
	return ds.RunInTransaction(ctx, func(ctx ae.Context) error {
		var dsErr error
		dsErr = f(ctx)
		return dsErr
	}, &ds.TransactionOptions{XG: crossGroup})
}
