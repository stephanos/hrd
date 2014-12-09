package internal

import (
	"github.com/101loops/hrd/internal/trafo"
	"github.com/101loops/hrd/internal/types"

	ae "appengine"
	ds "appengine/datastore"
)

// Count returns the number of results for a query.
func Count(ctx ae.Context, qry *types.Query) (int, error) {
	return qry.ToDSQuery(ctx).Count(ctx)
}

// Iterate loads entities from an iterator.
func Iterate(it *types.Iterator, dsts interface{}, multi bool) (keys []*types.Key, err error) {

	var docList *trafo.DocList
	if dsts != nil {
		docList, err = trafo.NewWriteableDocList(dsts, nil, multi)
		if err != nil {
			return
		}
	}

	var dsDocs []*trafo.Doc
	var dsKeys []*ds.Key
	for i := 0; ; i++ {

		var doc *trafo.Doc
		if docList != nil {
			doc, err = docList.Get(i)
			if err != nil {
				return
			}
			dsDocs = append(dsDocs, doc)
		}

		var dsKey *ds.Key
		dsKey, err = it.Next(doc.Pipe)
		if err == ds.Done {
			if !multi {
				if doc != nil {
					doc.Nil()
				}
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

	if docList == nil {
		return types.ImportKeys(dsKeys...), nil
	}

	keys, err = docList.ApplyResult(dsKeys, err)
	if dsDocs != nil {
		for i := range keys {
			docList.Add(keys[i], dsDocs[i])
		}
	}
	return keys, err
}
