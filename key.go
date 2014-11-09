package hrd

import (
	"fmt"
	"time"

	"appengine/datastore"
)

// Key represents the datastore key for an entity.
// It also contains meta data about said entity.
type Key struct {
	*datastore.Key

	// synced is the last time the entity was read/written.
	synced time.Time

	// err contains an error if the entity could not be loaded/saved.
	err *error
}

type keyBatch struct {
	keys   []*Key
	lo, hi int
}

// newKey creates a Key from a datastore.Key.
func newKey(k *datastore.Key) *Key {
	return &Key{Key: k}
}

// newKeys creates a sequence of Key from a sequence of datastore.Key.
func newKeys(keys []*datastore.Key) []*Key {
	ret := make([]*Key, len(keys))
	for i, k := range keys {
		ret[i] = newKey(k)
	}
	return ret
}

// Exists is whether an entity with this key exists in the datastore.
func (key *Key) Exists() bool {
	return !key.synced.IsZero()
}

// IDString returns the ID of this Key as a string.
func (key *Key) IDString() (id string) {
	id = key.StringID()
	if id == "" && key.IntID() > 0 {
		id = fmt.Sprintf("%v", key.IntID())
	}
	return
}

func (key *Key) String() string {
	return fmt.Sprintf("Key{'%v', %v}", key.Kind(), key.IDString())
}

func (key *Key) applyTo(src interface{}) {
	var parentKey = key.Parent()
	if parentKey != nil {
		id := parentKey.IntID()
		if parent, ok := src.(numParent); id != 0 && ok {
			parent.SetParent(id)
		} else {
			sid := parentKey.StringID()
			if parent, ok := src.(textParent); sid != "" && ok {
				parent.SetParent(sid)
			}
		}
	}

	id := key.IntID()
	if ident, ok := src.(numIdentifier); id != 0 && ok {
		ident.SetID(id)
	} else {
		sid := key.StringID()
		if ident, ok := src.(textIdentifier); sid != "" && ok {
			ident.SetID(sid)
		}
	}
}

// toDSKeys converts a sequence of Key to a sequence of datastore.Key.
func toDSKeys(keys []*Key) []*datastore.Key {
	ret := make([]*datastore.Key, len(keys))
	for i, k := range keys {
		ret[i] = k.Key
	}
	return ret
}
