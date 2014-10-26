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
	keyBatches := toKeyBatches(keys, putMultiLimit)
	for _, keyBatch := range keyBatches {
		docsBatch := docs.list[keyBatch.lo:keyBatch.hi]
		dsKeys, err := datastore.PutMulti(store.ctx, toDSKeys(keyBatch.keys), docsBatch)
		// TODO: appengine.MultiError
		if err != nil {
			return nil, store.logErr(err)
		}

		now := time.Now()
		for i := range keyBatch.keys {
			keyIdx := keyBatch.lo + i
			doc := docs.list[keyIdx]

			if keys[keyIdx].Incomplete() {
				keys[keyIdx] = newKey(dsKeys[keyIdx])
				doc.setKey(keys[keyIdx])
			}

			keys[keyIdx].opts = opts
			keys[keyIdx].synced = now

			toCache[keys[keyIdx]] = doc
		}
	}

	// #4 update cache
	store.cache.write(toCache)

	return keys, nil
}
