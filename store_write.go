package hrd

import (
	"appengine/datastore"
	"fmt"
	"time"
)

const putMultiLimit = 500
const deleteMultiLimit = 500

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
		for i, k := range keys {
			if k.Incomplete() {
				return nil, fmt.Errorf("incomplete key %q for %q (%dth index)", k, kind, i)
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
		if err != nil {
			return nil, store.logErr(err)
		}

		now := time.Now()
		for i, _ := range keys[lo:hi] {
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

func (store *Store) deleteMulti(kind string, keys []*Key) error {

	store.ctx.Infof(store.logAct("deleting", "from", keys, kind))

	// #1 delete from cache
	defer store.cache.delete(keys)

	// #2 delete from datastore
	for i := 0; i <= len(keys)/deleteMultiLimit; i++ {
		lo := i * deleteMultiLimit
		hi := (i + 1) * deleteMultiLimit
		if hi > len(keys) {
			hi = len(keys)
		}
		if err := datastore.DeleteMulti(store.ctx, toDSKeys(keys[lo:hi])); err != nil {
			return err
		}
	}

	return nil
}
