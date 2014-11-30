package types

import (
	"fmt"
	"reflect"
	"time"

	"github.com/101loops/hrd/entity"

	ae "appengine"
	ds "appengine/datastore"
)

// DSKeyConverter can return a datastore.Key.
type DSKeyConverter interface {
	ToDSKey(ctx ae.Context) *ds.Key
}

// Key represents the identifier for an entity.
type Key struct {
	*KeyState

	Kind      string
	Namespace string

	StringID string
	IntID    int64

	Parent *Key
}

// KeyState represents meta data of the datastore operation
// the key originates from, if any.
type KeyState struct {

	// Synced is the last time the entity was read/written.
	Synced *time.Time

	// Error contains an error if the key could not be loaded/saved.
	Error error
}

func NewKey(kind, stringID string, intID int64, parent *Key) *Key {
	return &Key{
		Kind: kind, StringID: stringID, IntID: intID, Parent: parent, KeyState: &KeyState{},
	}
}

// ImportKey creates a Key from a datastore.Key.
func ImportKey(dsKey *ds.Key) *Key {
	if dsKey == nil {
		return nil
	}
	key := NewKey(dsKey.Kind(), dsKey.StringID(), dsKey.IntID(), ImportKey(dsKey.Parent()))
	key.Namespace = dsKey.Namespace()
	return key
}

// ImportKeys creates a sequence of Key from a sequence of datastore.Key.
func ImportKeys(keys ...*ds.Key) []*Key {
	ret := make([]*Key, len(keys))
	for i, k := range keys {
		ret[i] = ImportKey(k)
	}
	return ret
}

// Incomplete returns whether the key does not refer to a stored entity.
// In particular, whether the key has a zero StringID and a zero IntID.
func (k *Key) Incomplete() bool {
	return k.StringID == "" && k.IntID == 0
}

func (k *Key) ToDSKey(ctx ae.Context) *ds.Key {
	var parentKey *ds.Key
	if k.Parent != nil {
		parentKey = k.Parent.ToDSKey(ctx)
	}
	return ds.NewKey(ctx, k.Kind, k.StringID, k.IntID, parentKey)
}

func (k *Key) String() string {
	return keyToString(k)
}

// keyToString returns a string representation of the passed-in datastore.Key.
func keyToString(key *Key) string {
	if key == nil {
		return ""
	}
	ret := fmt.Sprintf("Key{'%v', %v}", key.Kind, keyStringID(key))
	parent := keyToString(key.Parent)
	if parent != "" {
		ret += fmt.Sprintf("[Parent%v]", parent)
	}
	return ret
}

// keyStringID returns the ID of the passed-in datastore.Key as a string.
func keyStringID(key *Key) (id string) {
	id = key.StringID
	if id == "" && key.IntID > 0 {
		id = fmt.Sprintf("%v", key.IntID)
	}
	return
}

// GetEntityKey extracts a new Key from the given entity.
func GetEntityKey(kind *Kind, src interface{}) (*Key, error) {
	var parentKey *Key
	if parentIdent, ok := src.(entity.ParentNumIdentifier); ok {
		kind, id := parentIdent.Parent()
		parentKey = NewKey(kind, "", id, nil)
	}
	if parentIdent, ok := src.(entity.ParentTextIdentifier); ok {
		kind, id := parentIdent.Parent()
		parentKey = NewKey(kind, id, 0, nil)
	}

	if ident, ok := src.(entity.NumIdentifier); ok {
		return NewKey(kind.Name, "", ident.ID(), parentKey), nil
	}
	if ident, ok := src.(entity.TextIdentifier); ok {
		return NewKey(kind.Name, ident.ID(), 0, parentKey), nil
	}

	return nil, fmt.Errorf("value type %q does not provide ID()", reflect.TypeOf(src))
}

// GetEntitiesKeys extracts a sequence of Key from the given entities.
func GetEntitiesKeys(kind *Kind, src interface{}) ([]*Key, error) {
	srcVal := reflect.Indirect(reflect.ValueOf(src))
	srcKind := srcVal.Kind()
	if srcKind != reflect.Slice && srcKind != reflect.Map {
		return nil, fmt.Errorf("value must be a slice or map, but is %q", srcKind)
	}

	collLen := srcVal.Len()
	keys := make([]*Key, collLen)

	if srcVal.Kind() == reflect.Slice {
		for i := 0; i < collLen; i++ {
			v := srcVal.Index(i)
			key, err := GetEntityKey(kind, v.Interface())
			if err != nil {
				return nil, err
			}
			keys[i] = key
		}
		return keys, nil
	}

	for i, key := range srcVal.MapKeys() {
		v := srcVal.MapIndex(key)
		key, err := GetEntityKey(kind, v.Interface())
		if err != nil {
			return nil, err
		}
		keys[i] = key
	}
	return keys, nil
}
