package internal

import (
	"fmt"
	"time"

	"appengine"
	"appengine/datastore"
)

// LogDatastoreAction logs a datastore action.
func LogDatastoreAction(verb string, prop string, keys []*Key, kind string) string {
	if len(keys) == 1 {
		sKey := keys[0].String()
		sKey = "'" + sKey + "'"
		return fmt.Sprintf("%v %v", verb, sKey)
	}
	return fmt.Sprintf("%v %v items %v %q", verb, len(keys), prop, kind)
}

func applyResult(dsDocs []*doc, dsKeys []*datastore.Key, dsErr error) ([]*Key, error) {
	now := time.Now()
	keys := make([]*Key, len(dsKeys))

	var mErr appengine.MultiError
	if dsErr, ok := dsErr.(appengine.MultiError); ok {
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
