package hrd

import (
	"fmt"
	"reflect"
)

type docs struct {
	i        int
	list     []*doc
	keyList  []*Key
	srcVal   reflect.Value
	srcKind  reflect.Kind
	keyType  reflect.Type
	elemType reflect.Type
}

var (
	typeOfKey   = reflect.TypeOf(&Key{})
	typeOfInt64 = reflect.TypeOf(int64(0))
	typeOfStr   = reflect.TypeOf("")
)

func newReadableDocs(coll *Collection, src interface{}) (*docs, error) {
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

	// read collection
	srcColl := reflect.Indirect(srcVal)
	srcCollLen := srcColl.Len()

	// generate list of doc
	keys := make([]*Key, srcCollLen)
	list := make([]*doc, srcCollLen)
	for i := 0; i < srcCollLen; i++ {
		entity := srcColl.Index(i).Interface()

		d, err := newDocFromInst(entity)
		if err != nil {
			return nil, err
		}

		key, err := coll.getKey(entity)
		if err != nil {
			return nil, err
		}
		keys[i] = key

		list[i] = d
	}

	return &docs{list: list, keyList: keys}, nil
}

func newWriteableDocs(src interface{}, keys []*Key, multi bool) (*docs, error) {
	ret := &docs{keyList: keys}
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
		// create and validate collection
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

		// make sure the collection's elements are structs
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
			ret.add(key, d)
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

		ret.list = []*doc{d}
		ret.elemType = srcType
	}

	return ret, nil
}

func (docs *docs) keys() []*Key {
	return docs.keyList
}

func (docs *docs) set(idx int, src interface{}) {
	docs.list[idx].set(src)
}

func (docs *docs) get(idx int) *doc {
	return docs.list[idx]
}

func (docs *docs) nil(idx int) {
	docs.list[idx].nil()
}

func (docs *docs) next() (ret *doc, err error) {
	if docs.i < len(docs.list) {
		ret = docs.list[docs.i]
	} else {
		ret, err = newDocFromType(docs.elemType)
	}
	docs.i++
	return
}

func (docs *docs) add(key *Key, d *doc) {
	docs.list = append(docs.list, d)
	d.setKey(key)

	if docs.srcKind == reflect.Map {
		var v reflect.Value
		switch docs.keyType {
		case typeOfInt64:
			v = reflect.ValueOf(key.IntID())
		case typeOfStr:
			v = reflect.ValueOf(key.StringID())
		default:
			v = reflect.ValueOf(key)
		}
		docs.srcVal.SetMapIndex(v, d.srcVal)
	} else if docs.srcKind == reflect.Slice {
		docs.srcVal.Set(reflect.Append(docs.srcVal, d.srcVal))
	}
}
