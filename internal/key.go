package internal

import (
	"fmt"
	"reflect"
	"time"

	"github.com/101loops/hrd/entity"

	"appengine/datastore"
)

// Key represents the datastore key for an
// It also contains meta data about said
type Key struct {
	*datastore.Key

	// synced is the last time the entity was read/written.
	synced time.Time

	// err contains an error if the entity could not be loaded/saved.
	err *error
}

// NewKey creates a Key from a datastore.Key.
func NewKey(k *datastore.Key) *Key {
	return &Key{Key: k}
}

// newKeys creates a sequence of Key from a sequence of datastore.Key.
func newKeys(keys []*datastore.Key) []*Key {
	ret := make([]*Key, len(keys))
	for i, k := range keys {
		ret[i] = NewKey(k)
	}
	return ret
}

// Exists is whether an entity with this key exists in the datastore.
func (key *Key) Exists() bool {
	return !key.synced.IsZero()
}

// Error returns an error associated with the key.
func (key *Key) Error() error {
	if key.err == nil {
		return nil
	}
	return *key.err
}

func (key *Key) String() string {
	return keyToString(key.Key)
}

// keyStringID returns the ID of the passed-in Key as a string.
func keyStringID(key *datastore.Key) (id string) {
	id = key.StringID()
	if id == "" && key.IntID() > 0 {
		id = fmt.Sprintf("%v", key.IntID())
	}
	return
}

// keyToString returns a string representation of the passed-in Key.
func keyToString(key *datastore.Key) string {
	if key == nil {
		return ""
	}
	ret := fmt.Sprintf("Key{'%v', %v}", key.Kind(), keyStringID(key))
	parent := keyToString(key.Parent())
	if parent != "" {
		ret += fmt.Sprintf("[Parent%v]", parent)
	}
	return ret
}

// toDSKeys converts a sequence of Key to a sequence of datastore.Key.
func toDSKeys(keys []*Key) []*datastore.Key {
	ret := make([]*datastore.Key, len(keys))
	for i, k := range keys {
		ret[i] = k.Key
	}
	return ret
}

func getKey(kind Kind, src interface{}) (*Key, error) {
	ctx := kind.Context()

	var parentKey *datastore.Key
	if parented, ok := src.(entity.NumParent); ok {
		parentKey = datastore.NewKey(ctx, parented.ParentKind(), "", parented.Parent(), nil)
	}
	if parented, ok := src.(entity.TextParent); ok {
		parentKey = datastore.NewKey(ctx, parented.ParentKind(), parented.Parent(), 0, nil)
	}

	if ident, ok := src.(entity.NumIdentifier); ok {
		return NewKey(datastore.NewKey(ctx, kind.Name(), "", ident.ID(), parentKey)), nil
	}
	if ident, ok := src.(entity.TextIdentifier); ok {
		return NewKey(datastore.NewKey(ctx, kind.Name(), ident.ID(), 0, parentKey)), nil
	}
	return nil, fmt.Errorf("value type %q does not provide ID()", reflect.TypeOf(src))
}

func getKeys(kind Kind, src interface{}) ([]*Key, error) {
	srcVal := reflect.Indirect(reflect.ValueOf(src))
	srcKind := srcVal.Kind()
	if srcKind != reflect.Slice && srcKind != reflect.Map {
		return nil, fmt.Errorf("value must be a slice or map")
	}

	collLen := srcVal.Len()
	keys := make([]*Key, collLen)

	if srcVal.Kind() == reflect.Slice {
		for i := 0; i < collLen; i++ {
			v := srcVal.Index(i)
			key, err := getKey(kind, v.Interface())
			if err != nil {
				return nil, err
			}
			keys[i] = key
		}
		return keys, nil
	}

	for i, key := range srcVal.MapKeys() {
		v := srcVal.MapIndex(key)
		key, err := getKey(kind, v.Interface())
		if err != nil {
			return nil, err
		}
		keys[i] = key
	}
	return keys, nil
}
