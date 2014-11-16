package internal

import (
	"fmt"
	"time"

	ae "appengine"
	ds "appengine/datastore"
)

// LogDatastoreAction logs a datastore action.
func LogDatastoreAction(verb string, prop string, keys []*Key, kind string) string {
	if len(keys) == 1 {
		sKey := keys[0].String()
		return fmt.Sprintf("%v %v", verb, sKey)
	}
	return fmt.Sprintf("%v %v items %v %q", verb, len(keys), prop, kind)
}

func applyResult(dsDocs []*doc, dsKeys []*ds.Key, dsErr error) ([]*Key, error) {
	now := time.Now()
	keys := make([]*Key, len(dsKeys))

	var mErr ae.MultiError
	if dsErr, ok := dsErr.(ae.MultiError); ok {
		mErr = dsErr
	}

	hasErr := false
	for i := range dsKeys {
		keys[i] = NewKey(dsKeys[i])

		if mErr == nil || mErr[i] == nil {
			if dsDocs != nil {
				dsDocs[i].setKey(keys[i])
			}
			keys[i].synced = now
			continue
		}

		if mErr[i] == ds.ErrNoSuchEntity {
			dsDocs[i].nil() // not found: set to 'nil'
			mErr[i] = nil   // ignore error
			continue
		}

		hasErr = true
		keys[i].err = mErr[i]
	}

	if hasErr {
		return keys, mErr
	}
	return keys, nil
}
