package internal

import (
	"github.com/101loops/hrd/internal/trafo"
	"github.com/101loops/hrd/internal/types"

	ds "appengine/datastore"
)

// DSIterate loads entities from an iterator.
func DSIterate(dsIt *ds.Iterator, dsts interface{}, multi bool) (keys []*types.Key, err error) {

	// in a keys-only query there is no dsts
	var docs *trafo.Docs
	if dsts != nil {
		docs, err = trafo.NewWriteableDocs(dsts, nil, multi)
		if err != nil {
			return
		}
	}

	var dsDocs []*trafo.Doc
	var dsKeys []*ds.Key
	for {
		var doc *trafo.Doc
		if docs != nil {
			doc, err = docs.Next()
			if err != nil {
				return
			}
			dsDocs = append(dsDocs, doc)
		}

		var dsKey *ds.Key
		dsKey, err = dsIt.Next(doc)
		if err == ds.Done {
			if !multi {
				docs.Nil(0)
				return nil, nil
			}
			break
		}
		if err != nil {
			return
		}

		dsKeys = append(dsKeys, dsKey)

		if !multi {
			break
		}
	}

	keys, err = applyResult(dsDocs, dsKeys, err)
	if dsDocs != nil {
		for i := range keys {
			docs.Add(keys[i], dsDocs[i])
		}
	}

	return keys, err
}
