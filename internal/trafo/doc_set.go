package trafo

import (
	"fmt"
	"reflect"

	"github.com/101loops/hrd/internal/types"
)

// DocSet represents a collection of Doc.
type DocSet struct {
	list     []*Doc
	keyList  []*types.Key
	srcVal   reflect.Value
	srcKind  reflect.Kind
	keyType  reflect.Type
	elemType reflect.Type
}

var (
	typeOfKey   = reflect.TypeOf(&types.Key{})
	typeOfInt64 = reflect.TypeOf(int64(0))
)

// NewReadableDocSet returns a new DocSet, suitable for reading from it.
func NewReadableDocSet(kind *types.Kind, src interface{}) (*DocSet, error) {
	srcVal := reflect.ValueOf(src)
	if src == nil || (srcVal.Kind() == reflect.Ptr && srcVal.IsNil()) {
		return nil, fmt.Errorf("value must be non-nil")
	}

	srcKind := srcVal.Kind()
	switch srcKind {
	case reflect.Slice, reflect.Map:
	default:
		// wrap single entity in a slice
		src = []interface{}{src}
		srcVal = reflect.ValueOf(src)
	}

	// read Kind
	srcColl := reflect.Indirect(srcVal)
	srcCollLen := srcColl.Len()

	// generate list of doc
	keys := make([]*types.Key, srcCollLen)
	list := make([]*Doc, srcCollLen)
	for i := 0; i < srcCollLen; i++ {
		entity := srcColl.Index(i).Interface()

		d, err := newDocFromInst(entity)
		if err != nil {
			return nil, err
		}

		key, err := types.GetEntityKey(kind, entity)
		if err != nil {
			return nil, err
		}
		keys[i] = key

		list[i] = d
	}

	return &DocSet{list: list, keyList: keys}, nil
}

// NewWriteableDocSet creates a new DocSet, suitable for writing to it.
func NewWriteableDocSet(src interface{}, keys []*types.Key, multi bool) (*DocSet, error) {
	ret := &DocSet{keyList: keys}
	keysLen := len(keys)

	// resolve pointer
	srcVal := reflect.ValueOf(src)
	srcKind := srcVal.Kind()
	if srcKind != reflect.Ptr || src == nil || srcVal.IsNil() {
		return nil, fmt.Errorf("invalid value kind %q (wanted non-nil pointer)", srcKind)
	}
	srcVal = srcVal.Elem()
	srcType := srcVal.Type()

	if multi {
		// create and validate Kind
		srcKind = srcVal.Kind()
		switch srcKind {
		case reflect.Slice:
			// create new slice
			srcVal.Set(reflect.MakeSlice(srcType, 0, keysLen))
		case reflect.Map:
			// validate type of map key
			ret.keyType = srcType.Key()
			switch ret.keyType {
			case typeOfInt64, typeOfStr, typeOfKey:
			default:
				return nil, fmt.Errorf("invalid value key %q (wanted string, int64 or *hrd.Key)", srcKind)
			}
			// create new map
			srcVal.Set(reflect.MakeMap(srcType))
		default:
			return nil, fmt.Errorf("invalid value kind %q (wanted map or slice)", srcKind)
		}
		ret.srcKind = srcKind

		// make sure the Kind's elements are structs
		ret.elemType = srcType.Elem()
		collElemKind := ret.elemType.Kind()
		switch collElemKind {
		case reflect.Struct:
		case reflect.Ptr:
			collElemKind := ret.elemType.Elem().Kind()
			if collElemKind != reflect.Struct {
				return nil, fmt.Errorf("invalid value element kind (%q)", collElemKind)
			}
		default:
			return nil, fmt.Errorf("invalid value element kind (%q)", collElemKind)
		}

		// generate list of doc
		ret.srcVal = srcVal
		for _, key := range keys {
			d, err := newDocFromType(ret.elemType)
			if err != nil {
				return nil, err
			}
			ret.Add(key, d)
		}
	} else {
		if keysLen > 1 {
			return nil, fmt.Errorf("wanted exactly 1 key (got %d)", keysLen)
		}

		srcVal.Set(reflect.New(srcType.Elem()))
		d, err := newDoc(srcVal)
		if err != nil {
			return nil, err
		}

		ret.list = []*Doc{d}
		ret.elemType = srcType
	}

	return ret, nil
}

// List returns the set's sequence of Doc.
func (set *DocSet) List() []*Doc {
	return set.list
}

// Keys returns the set's sequence of Key.
func (set *DocSet) Keys() []*types.Key {
	return set.keyList
}

// Get returns the set's nth Doc.
// it is created first if it doesn't already exist.
func (set *DocSet) Get(nth int) (ret *Doc, err error) {
	if nth < len(set.list) {
		ret = set.list[nth]
	} else {
		ret, err = newDocFromType(set.elemType)
	}
	return
}

// Add appends a new Doc to the set.
func (set *DocSet) Add(key *types.Key, doc *Doc) {
	set.list = append(set.list, doc)
	doc.SetKey(key)

	if set.srcKind == reflect.Map {
		var v reflect.Value
		switch set.keyType {
		case typeOfInt64:
			v = reflect.ValueOf(key.IntID())
		case typeOfStr:
			v = reflect.ValueOf(key.StringID())
		default:
			v = reflect.ValueOf(key)
		}
		set.srcVal.SetMapIndex(v, doc.srcVal)
	} else if set.srcKind == reflect.Slice {
		set.srcVal.Set(reflect.Append(set.srcVal, doc.srcVal))
	}
}
