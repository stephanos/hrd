package internal

import (
	"fmt"

	"github.com/101loops/hrd/internal/trafo"
	"github.com/101loops/hrd/internal/types"
	"github.com/qedus/nds"

	ae "appengine"
	ds "appengine/datastore"
)

var (
	dsGet = func(ctx ae.Context, keys []*ds.Key, dst interface{}) error {
		return ds.GetMulti(ctx, keys, dst)
	}

	ndsGet = func(ctx ae.Context, keys []*ds.Key, dst interface{}) error {
		return nds.GetMulti(ctx, keys, dst)
	}
)

// Get loads entities for the given keys.
func Get(kind *types.Kind, keys []*types.Key, dst interface{}, useGlobalCache bool, multi bool) ([]*types.Key, error) {
	if err := validateGetKeys(kind, keys); err != nil {
		return nil, err
	}

	docs, err := trafo.NewWriteableDocList(dst, keys, multi)
	if err != nil {
		return nil, err
	}
	docsPipe := docs.Pipe(kind.Context)

	ctx := kind.Context
	ctx.Infof(LogDatastoreAction("getting", "from", keys, kind.Name))

	var dsErr error
	dsKeys := toDSKeys(keys)
	if useGlobalCache {
		dsErr = ndsGet(ctx, dsKeys, docsPipe.Properties())
	}
	dsErr = dsGet(ctx, dsKeys, docsPipe.Properties())

	return applyResult(docsPipe.Docs, dsKeys, dsErr)
}

func validateGetKeys(kind *types.Kind, keys []*types.Key) error {
	if keys == nil || len(keys) == 0 {
		return fmt.Errorf("no keys provided")
	}

	for i, key := range keys {
		if key.Incomplete() {
			return fmt.Errorf("'%v' is incomplete (%dth index)", key, i)
		}
	}

	for _, k := range keys {
		keyKind := k.Kind()
		if keyKind != kind.Name {
			err := fmt.Errorf("invalid key kind '%v' for Kind '%v'", keyKind, kind.Name)
			return logErr(kind.Context, err)
		}
	}

	return nil
}
