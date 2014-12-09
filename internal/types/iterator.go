package types

import (
	ae "appengine"
	ds "appengine/datastore"
)

// Iterator is the result of running a query.
type Iterator struct {
	inner *ds.Iterator
	ctx   ae.Context
	query *Query
}

// NewIterator returns a new Iterator by executing the passed-in query.
func NewIterator(ctx ae.Context, query *Query) *Iterator {
	return &Iterator{
		inner: query.ToDSQuery(ctx).Run(ctx),
		query: query,
		ctx:   ctx,
	}
}

// Cursor returns a cursor for the Iterator's current location.
func (it *Iterator) Cursor() (string, error) {
	c, err := it.inner.Cursor()
	if err != nil {
		return "", err
	}
	return c.String(), nil
}

// Next returns the key of the next result. If the query is not keys-only,
// it also loads the entity stored for that key into a PropertyLoadSaver.
func (it *Iterator) Next(pipeFunc func(ae.Context) ds.PropertyLoadSaver) (*ds.Key, error) {
	pipe := pipeFunc(it.ctx)
	return it.inner.Next(pipe)
}
