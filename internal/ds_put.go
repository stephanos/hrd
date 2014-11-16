package internal

import (
	"fmt"

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
func DSPut(kind Kind, src interface{}, completeKeys bool) ([]*Key, error) {
	ctx := kind.Context()

	docs, err := newReadableDocs(kind, src)
	if err != nil {
		return nil, err
	}

	keys := docs.keyList
	if err := validateDSPutKeys(kind, keys, completeKeys); err != nil {
		return nil, err
	}

	ctx.Infof(LogDatastoreAction("putting", "in", keys, kind.Name()))

	dsDocs := docs.list
	dsKeys, dsErr := dsPut(ctx, toDSKeys(keys), dsDocs)
	if dsErr != nil {
		return nil, dsErr
	}

	return applyResult(dsDocs, dsKeys, dsErr)
}

func validateDSPutKeys(kind Kind, keys []*Key, completeKeys bool) error {
	if len(keys) == 0 {
		return fmt.Errorf("no keys provided for %q", kind.Name())
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
