package hrd

import (
	"appengine/datastore"
	"fmt"
	"time"
)

type Key struct {
	*datastore.Key
	err     *error
	source  string
	version int64
	synced  time.Time
	opts    *operationOpts
}

func newKey(k *datastore.Key) *Key {
	return &Key{Key: k}
}

func newKeys(keys []*datastore.Key) []*Key {
	ret := make([]*Key, len(keys))
	for i, k := range keys {
		ret[i] = newKey(k)
	}
	return ret
}

func (key *Key) Exists() bool {
	return !key.synced.IsZero()
}

func (key *Key) IdString() (id string) {
	id = key.StringID()
	if id == "" && key.IntID() > 0 {
		id = fmt.Sprintf("%v", key.IntID())
	}
	return
}

func setKey(src interface{}, key *Key) {
	var parentKey = key.Parent()
	if parentKey != nil {
		id := parentKey.IntID()
		if parent, ok := src.(numParent); id != 0 && ok {
			parent.SetParent(id)
		}
		sid := parentKey.StringID()
		if parent, ok := src.(textParent); sid != "" && ok {
			parent.SetParent(sid)
		}
	}

	id := key.IntID()
	if ident, ok := src.(numIdentifier); id != 0 && ok {
		ident.SetID(id)
	}
	sid := key.StringID()
	if ident, ok := src.(textIdentifier); sid != "" && ok {
		ident.SetID(sid)
	}

	if v, ok := src.(versioned); ok {
		key.version = v.Version()
	}
}

// Convert a Key to string.
func toMemKey(k *Key) string {
	return fmt.Sprintf("%v-%v", k.Encode(), k.version)
}

// Convert a slice of Keys to a slice of strings.
func toMemKeys(keys []*Key) (memKeys []string) {
	for _, k := range keys {
		if !k.Incomplete() {
			memKeys = append(memKeys, toMemKey(k))
		}
	}
	return
}

// Convert a slice of Keys to a slice of datastore Keys.
func toDSKeys(keys []*Key) []*datastore.Key {
	ret := make([]*datastore.Key, len(keys))
	for i, k := range keys {
		ret[i] = k.Key
	}
	return ret
}
