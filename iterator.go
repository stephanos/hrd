package hrd

import "github.com/101loops/hrd/internal/types"

// Iterator is the result of running a query.
type Iterator struct {
	inner *types.Iterator
}

func newIterator(qry *Query) *Iterator {
	return &Iterator{types.NewIterator(qry.ctx, qry.inner)}
}

// Cursor returns a cursor for the Iterator's current location.
func (it *Iterator) Cursor() (string, error) {
	return it.inner.Cursor()
}

// GetOne loads an entity from the Iterator into the passed destination.
func (it *Iterator) GetOne(dst interface{}) (*Key, error) {
	keys, err := it.get(dst, false)
	if len(keys) == 0 {
		return nil, err
	}
	return keys[0], err
}

// GetAll loads all entities from the Iterator into the passed destination.
func (it *Iterator) GetAll(dsts interface{}) ([]*Key, error) {
	return it.get(dsts, true)
}

func (it *Iterator) get(dsts interface{}, multi bool) ([]*Key, error) {
	keys, err := dsIterate(it.inner, dsts, multi)
	return importKeys(keys), err
}
