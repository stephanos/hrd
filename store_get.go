package hrd

import (
	"fmt"
	"time"

	"appengine"
	"appengine/datastore"
)

const getMultiLimit = 1000

func (store *Store) getMulti(kind string, docs *docs, opts *operationOpts) ([]*Key, error) {
	meta, keys, err := store.getMultiStats(kind, docs, opts)
	if err == nil {
		store.ctx.Infof(meta.string())
	} else {
		store.ctx.Errorf("%v: %v", meta.descr, err)
	}
	return keys, err
}

func (store *Store) getMultiStats(kind string, docs *docs, opts *operationOpts) (*meta, []*Key, error) {

	meta := &meta{}

	// #1 find entity keys
	keys := docs.keyList
	if len(keys) == 0 {
		return meta, nil, fmt.Errorf("no keys provided")
	}

	for i, key := range keys {
		if key.Incomplete() {
			return meta, nil, fmt.Errorf("incomplete key '%v' (%dth index)", key, i)
		}
		key.opts = opts
	}

	meta.descr = store.logAct("getting", "from", keys, kind)

	// #2 read from cache
	dsKeys, dsDocs := store.cache.read(keys, docs)
	for _, key := range keys {
		if key.source == sourceMemcache {
			meta.fromGlobalCache++
		} else if key.source == sourceMemory {
			meta.fromLocalCache++
		}
	}

	// #3 load from datastore
	docsToCache := make(map[*Key]*doc, 0)
	for i := 0; i <= len(dsKeys)/getMultiLimit; i++ {
		lo := i * getMultiLimit
		hi := (i + 1) * getMultiLimit
		if hi > len(dsKeys) {
			hi = len(dsKeys)
		}

		dsErr := datastore.GetMulti(store.ctx, toDSKeys(dsKeys[lo:hi]), dsDocs[lo:hi])
		var merr appengine.MultiError
		if dsErr != nil {
			if multi, ok := dsErr.(appengine.MultiError); ok {
				merr = multi
			} else {
				return meta, keys, dsErr
			}
		}

		now := time.Now()
		for i, key := range dsKeys[lo:hi] {
			if merr == nil || merr[i] == nil {
				docsToCache[key] = dsDocs[lo+i]
				dsDocs[lo+i].setKey(key)
				key.source = sourceDatastore
				key.synced = now
				continue
			}

			if merr[i] == datastore.ErrNoSuchEntity {
				dsDocs[lo+i].nil() // not found: set to 'nil'
				merr[i] = nil      // ignore error
				continue
			}

			key.err = &merr[i]
		}
	}

	// #4 update cache
	store.cache.write(docsToCache)

	return meta, keys, nil
}
