package internal

import (
	"fmt"

	"github.com/qedus/nds"

	"appengine"
	"appengine/datastore"
)

// DSDelete deletes the entities for the given keys.
func DSDelete(ctx appengine.Context, kind string, keys []*datastore.Key) error {
	ctx.Infof(LogDatastoreAction("deleting", "from", keys, kind))
	return nds.DeleteMulti(ctx, keys)
}

// LogDatastoreAction logs a datastore action.
func LogDatastoreAction(verb string, prop string, keys []*datastore.Key, kind string) string {
	if len(keys) == 1 {
		id := KeyString(keys[0])
		if id == "" {
			return fmt.Sprintf("%v %v %v %q", verb, "1 item", prop, kind)
		}
		id = "'" + id + "'"

		parent := ""
		if parentKey := keys[0].Parent(); parentKey != nil {
			parent = fmt.Sprintf(" (with parent '%v')", KeyString(parentKey))
		}
		return fmt.Sprintf("%v %v %v %q%v", verb, id, prop, kind, parent)
	}

	return fmt.Sprintf("%v %v items %v %q", verb, len(keys), prop, kind)
}
