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
	ndsPut = func(ctx ae.Context, keys []*ds.Key, dst interface{}) ([]*ds.Key, error) {
		return nds.PutMulti(ctx, keys, dst)
	}
)

// Put saves the given entities.
func Put(kind *types.Kind, src interface{}, completeKeys bool) ([]*types.Key, error) {
	ctx := kind.Context

	docList, err := trafo.NewReadableDocList(kind, src)
	if err != nil {
		return nil, err
	}

	keys := docList.Keys()
	if err := validatePutKeys(kind, keys, completeKeys); err != nil {
		return nil, err
	}

	ctx.Infof(LogDatastoreAction("putting", "in", keys, kind.Name))

	docsPipe := docList.Pipe(kind.Context)
	dsKeys, dsErr := ndsPut(ctx, toDSKeys(ctx, keys), docsPipe.Properties())
	if dsErr != nil {
		return nil, dsErr
	}

	return docList.ApplyResult(dsKeys, dsErr)
}

func validatePutKeys(kind *types.Kind, keys []*types.Key, completeKeys bool) error {
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
