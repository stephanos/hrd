package internal

import (
	"fmt"
	"time"

	"github.com/101loops/hrd/entity"
	"github.com/qedus/nds"
)

// DSPut saves the given entities.
func DSPut(kind Kind, src interface{}, completeKeys bool) ([]*Key, error) {
	ctx := kind.Context()

	docs, err := newReadableDocs(kind, src)
	if err != nil {
		return nil, err
	}

	// get document keys
	keys := docs.keyList
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided for %q", kind.Name())
	}

	if completeKeys {
		for i, key := range keys {
			if key.Incomplete() {
				return nil, fmt.Errorf("%v is incomplete (%dth index)", key, i)
			}
		}
	}

	// timestamp documents
	for _, d := range docs.list {
		src := d.get()
		now := time.Now()
		if ts, ok := src.(entity.CreateTimestamper); ok {
			ts.SetCreatedAt(now)
		}
		if ts, ok := src.(entity.UpdateTimestamper); ok {
			ts.SetUpdatedAt(now)
		}
	}

	// put into datastore
	dsDocs := docs.list
	ctx.Infof(LogDatastoreAction("putting", "in", keys, kind.Name()))
	dsKeys, dsErr := nds.PutMulti(ctx, toDSKeys(keys), dsDocs)

	if dsErr != nil {
		return nil, dsErr
	}

	return applyResult(dsDocs, dsKeys, dsErr)
}
