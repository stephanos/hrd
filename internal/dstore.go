package internal

import (
	"fmt"

	"github.com/qedus/nds"

	"appengine"
	"appengine/datastore"
)

// DSGet loads entities for the given keys.
func DSGet(ctx appengine.Context, kind string, keys []*datastore.Key, dst interface{}, useGlobalCache bool) error {
	ctx.Infof(LogDatastoreAction("getting", "from", keys, kind))
	if useGlobalCache {
		return nds.GetMulti(ctx, keys, dst)
	}
	return datastore.GetMulti(ctx, keys, dst)
}

// DSDelete deletes the entities for the given keys.
func DSDelete(ctx appengine.Context, kind string, keys []*datastore.Key) error {
	ctx.Infof(LogDatastoreAction("deleting", "from", keys, kind))
	return nds.DeleteMulti(ctx, keys)
}

// LogDatastoreAction logs a datastore action.
func LogDatastoreAction(verb string, prop string, keys []*datastore.Key, kind string) string {
	if len(keys) == 1 {
		sKey := KeyString(keys[0])
		sKey = "'" + sKey + "'"
		return fmt.Sprintf("%v %v", verb, sKey)
	}
	return fmt.Sprintf("%v %v items %v %q", verb, len(keys), prop, kind)
}
