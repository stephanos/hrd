package hrd

import (
	"appengine/datastore"
)

// Iterator is the result of running a query.
type Iterator struct {
	qry *Query
	it  *datastore.Iterator
}

// Cursor returns a cursor for the Iterator's current location.
func (it *Iterator) Cursor() (string, error) {
	if c, err := it.it.Cursor(); err == nil {
		return c.String(), nil
	}
	return "", nil
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

	var docs *docs
	if dsts != nil {
		docs, err = newWriteableDocs(dsts, nil, multi)
		if err != nil {
			return
		}
	}

	qryType := it.qry.typeOf
	store := it.qry.coll.store

	toCache := make(map[*Key]*doc, 0)
	for {
		var doc *doc
		if docs != nil {
			doc, err = docs.next()
			if err != nil {
				return
			}
		}

		// #1 load from iterator
		var dsKey *datastore.Key
		dsKey, err = it.it.Next(doc)
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
		key := newKey(dsKey)
		key.source = sourceDatastore
		key.opts = it.qry.opts
		keys = append(keys, key)

		if docs != nil {
			if qryType != projectQry && !store.tx {

				// #2 try to read entity from local cache
				fromCache := false
				if key.opts.readLocalCache {
					if cached, ok := store.getMemory(key); ok && cached != nil {
						key.source = sourceMemory
						fromCache = true
						doc.set(cached)
					}
				}

				if !fromCache {
					toCache[key] = doc
				}
			}

			docs.add(key, doc)
		}

		if !multi {
			break
		}
	}

	// #3 update cache
	it.qry.coll.store.cache.write(toCache)

	return keys, nil
}
