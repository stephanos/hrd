package types

import (
	"fmt"
	"reflect"
	"time"

	"github.com/101loops/hrd/entity"

	ds "appengine/datastore"
)

// Key represents the identifier for an entity.
type Key struct {
	*ds.Key
	*KeyResult
}

type KeyResult struct {
	// Synced is the last time the entity was read/written.
	Synced *time.Time

	// Error contains an error if the key could not be loaded/saved.
	Error error
}

// NewKey creates a Key from a datastore.Key.
func NewKey(k *ds.Key) *Key {
	return &Key{k, &KeyResult{}}
}

// newKeys creates a sequence of Key from a sequence of datastore.Key.
func NewKeys(keys ...*ds.Key) []*Key {
	ret := make([]*Key, len(keys))
	for i, k := range keys {
		ret[i] = NewKey(k)
	}
	return ret
}

func (key *Key) String() string {
	return keyToString(key.Key)
}

// keyToString returns a string representation of the passed-in datastore.Key.
func keyToString(key *ds.Key) string {
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

// keyStringID returns the ID of the passed-in datastore.Key as a string.
func keyStringID(key *ds.Key) (id string) {
	id = key.StringID()
	if id == "" && key.IntID() > 0 {
		id = fmt.Sprintf("%v", key.IntID())
	}
	return
}

func GetEntityKey(kind *Kind, src interface{}) (*Key, error) {
	ctx := kind.Context

	var parentKey *ds.Key
	if parentIdent, ok := src.(entity.ParentNumIdentifier); ok {
		kind, id := parentIdent.Parent()
		parentKey = ds.NewKey(ctx, kind, "", id, nil)
	}
	if parentIdent, ok := src.(entity.ParentTextIdentifier); ok {
		kind, id := parentIdent.Parent()
		parentKey = ds.NewKey(ctx, kind, id, 0, nil)
	}

	if ident, ok := src.(entity.NumIdentifier); ok {
		return NewKey(ds.NewKey(ctx, kind.Name, "", ident.ID(), parentKey)), nil
	}
	if ident, ok := src.(entity.TextIdentifier); ok {
		return NewKey(ds.NewKey(ctx, kind.Name, ident.ID(), 0, parentKey)), nil
	}

	return nil, fmt.Errorf("value type %q does not provide ID()", reflect.TypeOf(src))
}

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
