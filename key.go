package hrd

import (
	"github.com/101loops/hrd/internal/types"

	ae "appengine"
	ds "appengine/datastore"
)

// Key represents the datastore key of an entity.
type Key struct {
	state *types.KeyState

	kind      string
	namespace string

	stringID string
	intID    int64

	parent *Key
}

func newKey(key *types.Key) *Key {
	if key == nil {
		return nil
	}

	var parent *Key
	if parentKey := key.Parent(); parentKey != nil {
		parent = newKey(types.NewKey(parentKey))
	}

	return &Key{key.KeyState, key.Kind(), key.Namespace(), key.StringID(), key.IntID(), parent}
}

func newKeys(keys []*types.Key) []*Key {
	ret := make([]*Key, len(keys))
	for i, k := range keys {
		ret[i] = newKey(k)
	}
	return ret
}

func newNumKey(kind *Kind, id int64, parent *Key) *Key {
	return &Key{state: &types.KeyState{}, kind: kind.name, intID: id, parent: parent}
}

func newTextKey(kind *Kind, id string, parent *Key) *Key {
	return &Key{state: &types.KeyState{}, kind: kind.name, stringID: id, parent: parent}
}

// Kind returns the key's kind (also known as entity type).
func (k *Key) Kind() string {
	return k.kind
}

// StringID returns the key's string ID
// (also known as an entity name or key name), which may be "".
func (k *Key) StringID() string {
	return k.stringID
}

// IntID returns the key's integer ID, which may be 0.
func (k *Key) IntID() int64 {
	return k.intID
}

// Parent returns the key's parent key, which may be nil.
func (k *Key) Parent() *Key {
	return k.parent
}

// Namespace returns the key's namespace.
func (k *Key) Namespace() string {
	return k.namespace
}

// Exists is whether an entity with this key exists in the datastore.
func (k *Key) Exists() bool {
	if k.state == nil {
		return false
	}
	if t := k.state.Synced; t != nil {
		return !t.IsZero()
	}
	return false
}

// Error returns an error associated with the key.
func (k *Key) Error() error {
	if k.state == nil {
		return nil
	}
	return k.state.Error
}

// Incomplete returns whether the key does not refer to a stored entity.
// In particular, whether the key has a zero StringID and a zero IntID.
func (k *Key) Incomplete() bool {
	return k.stringID == "" && k.intID == 0
}

func (k *Key) toDSKey(ctx ae.Context) *ds.Key {
	var parentKey *ds.Key
	if k.parent != nil {
		parentKey = k.parent.toDSKey(ctx)
	}
	return ds.NewKey(ctx, k.kind, k.stringID, k.intID, parentKey)
}

func toInternalKeys(ctx ae.Context, keys []*Key) []*types.Key {
	ret := make([]*types.Key, len(keys))
	for i, k := range keys {
		ret[i] = types.NewKey(k.toDSKey(ctx))
	}
	return ret
}
