package hrd

import (
	"appengine/datastore"
	"fmt"
	"time"
)

// Key represents the datastore key for an entity.
// It also contains meta data about said entity.
type Key struct {
	*datastore.Key

	// source describes where the entity was read from.
	source  string

	// version is the entity's version.
	version int64

	// synced is the last time the entity was read/written.
	synced  time.Time

	// opts are the options to use for reading/writing the entity.
	opts    *operationOpts

	// err contains an error if the entity could not be loaded/saved.
	err     *error
}

// newKey creates a Key from a datastore Key.
func newKey(k *datastore.Key) *Key {
	return &Key{Key: k}
}

// newKeys creates a sequence of Key from a sequence of datastore Key.
func newKeys(keys []*datastore.Key) []*Key {
	ret := make([]*Key, len(keys))
	for i, k := range keys {
		ret[i] = newKey(k)
	}
	return ret
}

// Exists is whether the entity with this key exists in the datastore.
func (key *Key) Exists() bool {
	return !key.synced.IsZero()
}

// IdString returns the ID of this Key as a string.
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

// toMemKeys converts a sequence of Key to a sequence of string.
func toMemKeys(keys []*Key) []string {
	ret := make([]string, len(keys))
	for i, k := range keys {
		if !k.Incomplete() {
			ret[i] = toMemKey(k)
		}
	}
	return ret
}

// toDSKeys converts a sequence of Key to a sequence of datastore Key.
func toDSKeys(keys []*Key) []*datastore.Key {
	ret := make([]*datastore.Key, len(keys))
	for i, k := range keys {
		ret[i] = k.Key
	}
	return ret
}
