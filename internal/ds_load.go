package internal

import (
	"fmt"

	"github.com/qedus/nds"

	"appengine/datastore"
)

// DSGet loads entities for the given keys.
func DSGet(kind Kind, keys []*Key, dst interface{}, useGlobalCache bool, multi bool) ([]*Key, error) {
	ctx := kind.Context()
	collKind := kind.Name()

	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}

	for i, key := range keys {
		if key.Incomplete() {
			return nil, fmt.Errorf("'%v' is incomplete (%dth index)", key, i)
		}
	}

	for _, k := range keys {
		keyKind := k.Kind()
		if keyKind != collKind {
			err := fmt.Errorf("invalid key kind '%v' for Kind '%v'", keyKind, kind.Name())
			return nil, logErr(ctx, err)
		}
	}

	docs, err := newWriteableDocs(dst, keys, multi)
	if err != nil {
		return nil, err
	}

	dsDocs := docs.list
	dsKeys := toDSKeys(keys)
	ctx.Infof(LogDatastoreAction("getting", "from", keys, collKind))

	var dsErr error
	if useGlobalCache {
		dsErr = nds.GetMulti(ctx, dsKeys, dst)
	}
	dsErr = datastore.GetMulti(ctx, dsKeys, dst)
	return applyResult(dsDocs, dsKeys, dsErr)
}
