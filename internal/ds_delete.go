package internal

import (
	"github.com/qedus/nds"

	ae "appengine"
	ds "appengine/datastore"
)

var (
	dsDel = func(ctx ae.Context, keys []*ds.Key) error {
		return nds.DeleteMulti(ctx, keys)
	}
)

// DSDelete deletes the given entity.
func DSDelete(kind Kind, src interface{}, multi bool) error {
	var err error
	var keys []*Key

	if multi {
		keys, err = getKeys(kind, src)
	} else {
		var key *Key
		key, err = getKey(kind, src)
		keys = []*Key{key}
	}

	if err != nil {
		return err
	}

	return DSDeleteKeys(kind, keys)
}

// DSDeleteKeys deletes the entities for the given keys.
func DSDeleteKeys(kind Kind, keys []*Key) error {
	ctx := kind.Context()
	dsKeys := toDSKeys(keys)

	ctx.Infof(LogDatastoreAction("deleting", "from", keys, kind.Name()))

	return dsDel(ctx, dsKeys)
}
