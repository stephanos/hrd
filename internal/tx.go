package internal

import (
	ae "appengine"
	ds "appengine/datastore"
)

func DSTransact(ctx ae.Context, f func(_ ae.Context) error, crossGroup bool) error {
	return ds.RunInTransaction(ctx, func(ctx ae.Context) error {
		var dsErr error
		dsErr = f(ctx)
		return dsErr
	}, &ds.TransactionOptions{XG: crossGroup})
}
