package hrd

import (
	"appengine/datastore"
	"fmt"
	"time"
)

// Key represents the datastore key for a stored entity.
// It embeds the datastore's key; also it adds several internal
// fields.
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

// toMemKey converts a Key to a string. It includes the entity's version
// to prevent reading old versions of an entity from memcache.
func toMemKey(k *Key) string {
	return fmt.Sprintf("%v-%v", k.Encode(), k.version)
}

// toMemKeys converts a sequence of Keys to a sequence of strings.
func toMemKeys(keys []*Key) []string {
	ret := make([]string, len(keys))
	for i, k := range keys {
		if !k.Incomplete() {
			ret[i] = toMemKey(k)
		}
	}
	return ret
}

// toDSKeys converts a sequence of Keys to a sequence of datastore Keys.
func toDSKeys(keys []*Key) []*datastore.Key {
	ret := make([]*datastore.Key, len(keys))
	for i, k := range keys {
		ret[i] = k.Key
	}
	return ret
}
