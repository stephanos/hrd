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
	dsPut = func(ctx ae.Context, keys []*ds.Key, dst interface{}) ([]*ds.Key, error) {
		return nds.PutMulti(ctx, keys, dst)
	}
)

// DSPut saves the given entities.
func DSPut(kind *types.Kind, src interface{}, completeKeys bool) ([]*types.Key, error) {
	ctx := kind.Context

	docs, err := trafo.NewReadableDocSet(kind, src)
	if err != nil {
		return nil, err
	}

	keys := docs.Keys()
	if err := validateDSPutKeys(kind, keys, completeKeys); err != nil {
		return nil, err
	}

	ctx.Infof(LogDatastoreAction("putting", "in", keys, kind.Name))

	dsDocs := docs.List()
	dsKeys, dsErr := dsPut(ctx, toDSKeys(keys), dsDocs)
	if dsErr != nil {
		return nil, dsErr
	}

	return applyResult(dsDocs, dsKeys, dsErr)
}

func validateDSPutKeys(kind *types.Kind, keys []*types.Key, completeKeys bool) error {
	if len(keys) == 0 {
		return fmt.Errorf("no keys provided for %q", kind.Name)
	}

	if completeKeys {
		for i, key := range keys {
			if key.Incomplete() {
				return fmt.Errorf("%v is incomplete (%dth index)", key, i)
			}
		}
	}

	return nil
}
