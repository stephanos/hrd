package hrd

import (
	"github.com/101loops/hrd/internal/types"

	ae "appengine"
	ds "appengine/datastore"
)

// Key represents the datastore key of an entity.
type Key struct {
	inner *types.Key
}

func importKey(key *types.Key) *Key {
	if key == nil {
		return nil
	}
	return &Key{key}
}

func importKeys(keys []*types.Key) []*Key {
	ret := make([]*Key, len(keys))
	for i, key := range keys {
		ret[i] = importKey(key)
	}
	return ret
}

func newKey(kind string, stringID string, intID int64, parent *Key) *Key {
	var parentKey *types.Key
	if parent != nil {
		parentKey = parent.inner
	}
	return importKey(types.NewKey(kind, stringID, intID, parentKey))
}

func newNumKey(kind *Kind, id int64, parent *Key) *Key {
	return newKey(kind.name, "", id, parent)
}

func newTextKey(kind *Kind, id string, parent *Key) *Key {
	return newKey(kind.name, id, 0, parent)
}

// Kind returns the key's kind (also known as entity type).
func (k *Key) Kind() string {
	return k.inner.Kind
}

// StringID returns the key's string ID
// (also known as an entity name or key name), which may be empty.
func (k *Key) StringID() string {
	return k.inner.StringID
}

// IntID returns the key's integer ID, which may be zero.
func (k *Key) IntID() int64 {
	return k.inner.IntID
}

// Parent returns the key's parent key, which may be nil.
func (k *Key) Parent() *Key {
	return importKey(k.inner.Parent)
}

// Namespace returns the key's namespace.
func (k *Key) Namespace() string {
	return k.inner.Namespace
}

// Exists is whether an entity with this key exists in the datastore.
func (k *Key) Exists() bool {
	if t := k.inner.Synced; t != nil {
		return !t.IsZero()
	}
	return false
}

// Error returns an error associated with the key.
func (k *Key) Error() error {
	return k.inner.Error
}

// Incomplete returns whether the key does not refer to a stored entity.
// In particular, whether the key has a zero StringID and a zero IntID.
func (k *Key) Incomplete() bool {
	return k.inner.Incomplete()
}

// ToDSKey returns the respective datastore.Key.
func (k *Key) ToDSKey(ctx ae.Context) *ds.Key {
	return k.inner.ToDSKey(ctx)
}

func toInternalKeys(keys []*Key) []*types.Key {
	ret := make([]*types.Key, len(keys))
	for i, k := range keys {
		ret[i] = k.inner
	}
	return ret
}
