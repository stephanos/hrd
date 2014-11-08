package hrd

import "appengine/datastore"

// Iterator is the result of running a query.
type Iterator struct {
	qry  *Query
	dsIt *datastore.Iterator
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
func (it *Iterator) GetAll(dsts interface{}) (keys []*Key, err error) {
	return it.get(dsts, true)
}

func (it *Iterator) get(dsts interface{}, multi bool) (keys []*Key, err error) {
	if it.qry.err != nil {
		return nil, *it.qry.err
	}

	// in a keys-only query there is no dsts
	var docs *docs
	if dsts != nil {
		docs, err = newWriteableDocs(dsts, nil, multi)
		if err != nil {
			return
		}
	}

	var dsDocs []*doc
	var dsKeys []*datastore.Key
	for {
		var doc *doc
		if docs != nil {
			doc, err = docs.next()
			if err != nil {
				return
			}
			dsDocs = append(dsDocs, doc)
		}

		var dsKey *datastore.Key
		dsKey, err = it.dsIt.Next(doc)
		if err == datastore.Done {
			if !multi {
				docs.nil(0)
				return nil, nil
			}
			break
		}
		if err != nil {
			return
		}

		dsKeys = append(dsKeys, dsKey)

		if !multi {
			break
		}
	}

	keys, err = postProcess(dsDocs, dsKeys, err)
	if dsDocs != nil {
		for i := range keys {
			docs.add(keys[i], dsDocs[i])
		}
	}
	return keys, err
}
