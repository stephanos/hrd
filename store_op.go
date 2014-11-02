package hrd

import (
	"fmt"
	"time"

	"github.com/qedus/nds"

	"appengine"
	"appengine/datastore"
)

func (store *Store) getMulti(kind string, docs *docs, opts *operationOpts) ([]*Key, error) {

	keys := docs.keyList
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}

	for i, key := range keys {
		if key.Incomplete() {
			return nil, fmt.Errorf("incomplete key '%v' (%dth index)", key, i)
		}
	}

	store.ctx.Infof(store.logAct("getting", "from", keys, kind))

	var dsErr error
	dsDocs := docs.list
	dsKeys := toDSKeys(keys)
	if opts.useGlobalCache {
		dsErr = nds.GetMulti(store.ctx, dsKeys, dsDocs)
	}
	dsErr = datastore.GetMulti(store.ctx, dsKeys, dsDocs)

	return postProcess(dsDocs, dsKeys, dsErr)
}

func (store *Store) putMulti(kind string, docs *docs, opts *operationOpts) ([]*Key, error) {

	// get document keys
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

	// timestamp documents
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

	store.ctx.Infof(store.logAct("putting", "in", keys, kind))

	// put into datastore
	dsDocs := docs.list
	dsKeys, dsErr := nds.PutMulti(store.ctx, toDSKeys(keys), dsDocs)
	if dsErr != nil {
		return nil, store.logErr(dsErr)
	}

	return postProcess(dsDocs, dsKeys, dsErr)
}

func (store *Store) deleteMulti(kind string, keys []*Key) (err error) {
	store.ctx.Infof(store.logAct("deleting", "from", keys, kind))
	return nds.DeleteMulti(store.ctx, toDSKeys(keys))
}

func postProcess(dsDocs []*doc, dsKeys []*datastore.Key, dsErr error) ([]*Key, error) {
	now := time.Now()
	keys := make([]*Key, len(dsKeys))

	var mErr appengine.MultiError
	if dsErr, ok := dsErr.(appengine.MultiError); ok {
		mErr = dsErr
	}

	hasErr := false
	for i := range dsKeys {
		keys[i] = newKey(dsKeys[i])

		if mErr == nil || mErr[i] == nil {
			if dsDocs != nil {
				dsDocs[i].setKey(keys[i])
			}
			keys[i].synced = now
			continue
		}

		if mErr[i] == datastore.ErrNoSuchEntity {
			dsDocs[i].nil() // not found: set to 'nil'
			mErr[i] = nil   // ignore error
			continue
		}

		hasErr = true
		keys[i].err = &mErr[i]
	}

	if hasErr {
		return keys, mErr
	}
	return keys, nil
}
