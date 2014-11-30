package internal

import (
	"fmt"

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

func toDSKeys(ctx ae.Context, keys []*types.Key) []*ds.Key {
	ret := make([]*ds.Key, len(keys))
	for i, key := range keys {
		ret[i] = key.ToDSKey(ctx)
	}
	return ret
}
