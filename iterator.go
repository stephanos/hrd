package hrd

import (
	"appengine/datastore"
)

type Iterator struct {
	qry *Query
	it  *datastore.Iterator
}

// Cursor returns a cursor for the iterator's current location.
func (self *Iterator) Cursor() (string, error) {
	if c, err := self.it.Cursor(); err == nil {
		return c.String(), nil
	} else {
		return "", err
	}
}

func (self *Iterator) GetOne(dst interface{}) (err error) {
	_, err = self.get(dst, false)
	return
}

func (self *Iterator) GetAll(dsts interface{}) (keys []*Key, err error) {
	return self.get(dsts, true)
}

func (self *Iterator) get(dsts interface{}, multi bool) (keys []*Key, err error) {
	if self.qry.err != nil {
		return nil, *self.qry.err
	}

	var docs *docs
	if dsts != nil {
		docs, err = newWriteableDocs(dsts, nil, multi)
		if err != nil {
			return
		}
	}

	qryType := self.qry.typeOf
	store := self.qry.coll.store

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
		dsKey, err = self.it.Next(doc)
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
		key.source = SOURCE_DATASTORE
		key.opts = self.qry.opts
		keys = append(keys, key)

		if docs != nil {
			if qryType != ProjQry && !store.inTX {

				// #2 try to read entity from local cache
				fromCache := false
				if key.opts.readLocalCache {
					if cached, ok := store.getMemory(key); ok && cached != nil {
						key.source = SOURCE_MEMORY
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
	self.qry.coll.store.cache.write(toCache)

	return keys, nil
}
