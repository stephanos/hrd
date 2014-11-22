package internal

import (
	"fmt"
	"time"

	"github.com/101loops/hrd/internal/trafo"
	"github.com/101loops/hrd/internal/types"

	ae "appengine"
	ds "appengine/datastore"
)

// LogDatastoreAction logs a datastore action.
func LogDatastoreAction(verb string, prop string, keys []*types.Key, kind string) string {
	if len(keys) == 1 {
		sKey := keys[0].String()
		return fmt.Sprintf("%v %v", verb, sKey)
	}
	return fmt.Sprintf("%v %v items %v %q", verb, len(keys), prop, kind)
}

// toDSKeys converts a sequence of Key to a sequence of datastore.Key.
func toDSKeys(keys []*types.Key) []*ds.Key {
	ret := make([]*ds.Key, len(keys))
	for i, k := range keys {
		ret[i] = k.Key
	}
	return ret
}

func applyResult(dsDocs []*trafo.Doc, dsKeys []*ds.Key, dsErr error) ([]*types.Key, error) {
	now := time.Now()
	keys := make([]*types.Key, len(dsKeys))

	var mErr ae.MultiError
	if dsErr, ok := dsErr.(ae.MultiError); ok {
		mErr = dsErr
	}

	hasErr := false
	for i := range dsKeys {
		keys[i] = types.NewKey(dsKeys[i])

		if mErr == nil || mErr[i] == nil {
			if dsDocs != nil {
				dsDocs[i].SetKey(keys[i])
			}
			keys[i].Synced = now
			fmt.Printf("key: %v (%v)\n", keys[i], now)
			continue
		}

		if mErr[i] == ds.ErrNoSuchEntity {
			dsDocs[i].Nil() // not found: set to 'nil'
			mErr[i] = nil   // ignore error
			continue
		}

		hasErr = true
		keys[i].Error = mErr[i]
	}

	if hasErr {
		return keys, mErr
	}
	return keys, nil
}
