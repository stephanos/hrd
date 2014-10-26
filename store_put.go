package hrd

import (
	"fmt"
	"time"

	"appengine/datastore"
)

const putMultiLimit = 500

func (store *Store) putMulti(kind string, docs *docs, opts *operationOpts) ([]*Key, error) {

	// #1 timestamp documents
	for _, d := range docs.list {
		src := d.get()
		now := time.Now()
		if ts, ok := src.(timestampCreator); ok {
			ts.SetCreatedAt(now)
		}
		if ts, ok := src.(timestampUpdater); ok {
			ts.SetUpdatedAt(now)
		}
	}

	// #2 get document keys
	keys := docs.keyList
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided for %q", kind)
	}

	if opts.completeKeys {
		for i, key := range keys {
			if key.Incomplete() {
				return nil, fmt.Errorf("incomplete key %v for %q (%dth index)", key, kind, i)
			}
		}
	}

	store.ctx.Infof(store.logAct("putting", "in", keys, kind))

	// #3 put into datastore
	toCache := make(map[*Key]*doc, len(keys))
	for i := 0; i <= len(keys)/putMultiLimit; i++ {
		lo := i * putMultiLimit
		hi := (i + 1) * putMultiLimit
		if hi > len(keys) {
			hi = len(keys)
		}

		dsKeys, err := datastore.PutMulti(store.ctx, toDSKeys(keys[lo:hi]), docs.list[lo:hi])
		// TODO: appengine.MultiError
		if err != nil {
			return nil, store.logErr(err)
		}

		now := time.Now()
		for i := range keys[lo:hi] {
			doc := docs.list[lo+i]

			if keys[i].Incomplete() {
				keys[i] = newKey(dsKeys[i])
				doc.setKey(keys[i])
			}

			keys[i].opts = opts
			keys[i].synced = now

			toCache[keys[i]] = doc
		}
	}

	// #4 update cache
	store.cache.write(toCache)

	return keys, nil
}
