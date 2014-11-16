package internal

import ds "appengine/datastore"

// DSIterate loads entities from an iterator.
func DSIterate(dsIt *ds.Iterator, dsts interface{}, multi bool) (keys []*Key, err error) {

	// in a keys-only query there is no dsts
	var docs *docs
	if dsts != nil {
		docs, err = newWriteableDocs(dsts, nil, multi)
		if err != nil {
			return
		}
	}

	var dsDocs []*doc
	var dsKeys []*ds.Key
	for {
		var doc *doc
		if docs != nil {
			doc, err = docs.next()
			if err != nil {
				return
			}
			dsDocs = append(dsDocs, doc)
		}

		var dsKey *ds.Key
		dsKey, err = dsIt.Next(doc)
		if err == ds.Done {
			if !multi {
				docs.nil(0)
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
			docs.add(keys[i], dsDocs[i])
		}
	}

	return keys, err
}
