package hrd

import "appengine/datastore"

const deleteMultiLimit = 500

func (store *Store) deleteMulti(kind string, keys []*Key) (err error) {

	store.ctx.Infof(store.logAct("deleting", "from", keys, kind))

	// #1 delete from cache
	defer store.cache.delete(keys)

	// #2 delete from datastore
	keyBatches := toKeyBatches(keys, deleteMultiLimit)
	for _, keyBatch := range keyBatches {
		err = datastore.DeleteMulti(store.ctx, toDSKeys(keyBatch.keys))
		// TODO: appengine.MultiError
		if err != nil {
			return err
		}
	}

	return
}
