package hrd

import ds "appengine/datastore"

// Iterator is the result of running a query.
type Iterator struct {
	qry  *Query
	dsIt *ds.Iterator
}

// Cursor returns a cursor for the Iterator's current location.
func (it *Iterator) Cursor() (string, error) {
	c, err := it.dsIt.Cursor()
	if err != nil {
		return "", err
	}
	return c.String(), nil
}

// GetOne loads an entity from the Iterator into the passed destination.
func (it *Iterator) GetOne(dst interface{}) (err error) {
	_, err = it.get(dst, false)
	return
}

// GetAll loads all entities from the Iterator into the passed destination.
func (it *Iterator) GetAll(dsts interface{}) ([]*Key, error) {
	return it.get(dsts, true)
}

func (it *Iterator) get(dsts interface{}, multi bool) ([]*Key, error) {
	if it.qry.err != nil {
		return nil, *it.qry.err
	}

	keys, err := dsIterate(it.dsIt, dsts, multi)
	return newKeys(keys), err
}
