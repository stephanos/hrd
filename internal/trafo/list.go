package trafo

import (
	"fmt"
	"reflect"
	"time"

	"github.com/101loops/hrd/internal/types"

	ae "appengine"
	ds "appengine/datastore"
)

// DocList represents a collection of Doc.
type DocList struct {
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

// NewReadableDocList returns a new DocList, suitable for reading from it.
func NewReadableDocList(kind *types.Kind, src interface{}) (*DocList, error) {
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
		list[i] = d

		key, err := types.GetEntityKey(kind, entity)
		if err != nil {
			return nil, err
		}
		keys[i] = key
	}

	return &DocList{list: list, keyList: keys}, nil
}

// NewWriteableDocList creates a new DocList, suitable for writing to it.
func NewWriteableDocList(src interface{}, keys []*types.Key, multi bool) (*DocList, error) {
	ret := &DocList{keyList: keys}
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
		if collElemKind == reflect.Ptr {
			collElemKind = ret.elemType.Elem().Kind()
			if collElemKind != reflect.Struct {
				return nil, fmt.Errorf("invalid value element kind %q (wanted struct)", collElemKind)
			}
		} else {
			return nil, fmt.Errorf("invalid value element kind %q (wanted pointer)", collElemKind)
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

// Pipe returns the a DocsPipe to load/save the entities.
func (l *DocList) Pipe(ctx ae.Context) DocsPipe {
	return DocsPipe{ctx, l.list}
}

// Keys returns the list's sequence of Key.
func (l *DocList) Keys() []*types.Key {
	return l.keyList
}

// Get returns the list's nth Doc.
// it is created first if it doesn't already exist.
func (l *DocList) Get(nth int) (ret *Doc, err error) {
	if nth < len(l.list) {
		ret = l.list[nth]
	} else {
		ret, err = newDocFromType(l.elemType)
	}
	return
}

// Add appends a new Doc to the list.
func (l *DocList) Add(key *types.Key, doc *Doc) {
	l.list = append(l.list, doc)
	doc.setKey(key)

	if l.srcKind == reflect.Map {
		var v reflect.Value
		switch l.keyType {
		case typeOfInt64:
			v = reflect.ValueOf(key.IntID())
		case typeOfStr:
			v = reflect.ValueOf(key.StringID())
		default:
			v = reflect.ValueOf(key)
		}
		l.srcVal.SetMapIndex(v, doc.srcVal)
	} else if l.srcKind == reflect.Slice {
		l.srcVal.Set(reflect.Append(l.srcVal, doc.srcVal))
	}
}

func (l *DocList) ApplyResult(dsKeys []*ds.Key, dsErr error) ([]*types.Key, error) {
	now := time.Now()
	keys := make([]*types.Key, len(dsKeys))

	var mErr ae.MultiError
	if dsErr, ok := dsErr.(ae.MultiError); ok {
		mErr = dsErr
	}

	hasErr := false
	dsDocs := l.list
	for i := range dsKeys {
		keys[i] = types.NewKey(dsKeys[i])

		if mErr == nil || mErr[i] == nil {
			if dsDocs != nil {
				dsDocs[i].setKey(keys[i])
			}
			keys[i].Synced = &now
			continue
		}

		if mErr[i] == ds.ErrNoSuchEntity {
			dsDocs[i].Nil() // not found: set to 'nil'
			mErr[i] = nil   // ignore error
			continue
		}

		hasErr = true
		keys[i].Error = mErr[i]
	}

	if hasErr {
		return keys, mErr
	}
	return keys, nil
}
