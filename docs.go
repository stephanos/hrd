package hrd

import (
	"fmt"
	"reflect"
)

type docs struct {
	i         int
	list      []*doc
	keyList   []*Key
	src_val   reflect.Value
	src_kind  reflect.Kind
	key_type  reflect.Type
	elem_type reflect.Type
}

var (
	typeOfKey   = reflect.TypeOf(&Key{})
	typeOfInt64 = reflect.TypeOf(int64(0))
	typeOfStr   = reflect.TypeOf("")
)

func newReadableDocs(coll *Collection, src interface{}) (*docs, error) {
	src_val := reflect.ValueOf(src)
	if src == nil || (src_val.Kind() == reflect.Ptr && src_val.IsNil()) {
		return nil, fmt.Errorf("value must be non-nil")
	}

	src_kind := src_val.Kind()
	switch src_kind {
	case reflect.Slice, reflect.Map:
	default:
		// wrap single entity in a slice
		src = []interface{}{src}
		src_val = reflect.ValueOf(src)
	}

	// read collection
	src_coll := reflect.Indirect(src_val)
	src_coll_len := src_coll.Len()

	// generate list of doc
	keys := make([]*Key, src_coll_len)
	list := make([]*doc, src_coll_len)
	for i := 0; i < src_coll_len; i++ {
		entity := src_coll.Index(i).Interface()

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
	src_val := reflect.ValueOf(src)
	src_kind := src_val.Kind()
	if src_kind != reflect.Ptr || src == nil || src_val.IsNil() {
		return nil, fmt.Errorf("invalid value kind %q (wanted non-nil pointer)", src_kind)
	}
	src_val = src_val.Elem()
	src_type := src_val.Type()

	if multi {
		// create and validate collection
		src_kind = src_val.Kind()
		switch src_kind {
		case reflect.Slice:
			// create new slice
			src_val.Set(reflect.MakeSlice(src_type, 0, keysLen))
		case reflect.Map:
			// validate type of map key
			ret.key_type = src_type.Key()
			switch ret.key_type {
			case typeOfInt64, typeOfStr, typeOfKey:
			default:
				return nil, fmt.Errorf("invalid value key %q (wanted string, int64 or *hrd.Key)", src_kind)
			}
			// create new map
			src_val.Set(reflect.MakeMap(src_type))
		default:
			return nil, fmt.Errorf("invalid value kind %q (wanted map or slice)", src_kind)
		}
		ret.src_kind = src_kind

		// make sure the collection's elements are structs
		ret.elem_type = src_type.Elem()
		coll_elem_kind := ret.elem_type.Kind()
		switch coll_elem_kind {
		case reflect.Struct:
		case reflect.Ptr:
			coll_elem_kind := ret.elem_type.Elem().Kind()
			if coll_elem_kind != reflect.Struct {
				return nil, fmt.Errorf("invalid value element kind (%q)", coll_elem_kind)
			}
		default:
			return nil, fmt.Errorf("invalid value element kind (%q)", coll_elem_kind)
		}

		// generate list of doc
		ret.src_val = src_val
		for _, key := range keys {
			d, err := newDocFromType(ret.elem_type)
			if err != nil {
				return nil, err
			}
			ret.add(key, d)
		}
	} else {
		if keysLen > 1 {
			return nil, fmt.Errorf("wanted exactly 1 key (got %d)", keysLen)
		}

		src_val.Set(reflect.New(src_type.Elem()))
		d, err := newDoc(src_val)
		if err != nil {
			return nil, err
		}

		ret.list = []*doc{d}
		ret.elem_type = src_type
	}

	return ret, nil
}

func (self *docs) keys() []*Key {
	return self.keyList
}

func (self *docs) set(idx int, src interface{}) {
	self.list[idx].set(src)
}

func (self *docs) get(idx int) *doc {
	return self.list[idx]
}

func (self *docs) nil(idx int) {
	self.list[idx].nil()
}

func (self *docs) next() (ret *doc, err error) {
	if self.i < len(self.list) {
		ret = self.list[self.i]
	} else {
		ret, err = newDocFromType(self.elem_type)
	}
	self.i += 1
	return
}

func (self *docs) add(key *Key, d *doc) {
	self.list = append(self.list, d)
	d.setKey(key)

	if self.src_kind == reflect.Map {
		var v reflect.Value
		switch self.key_type {
		case typeOfInt64:
			v = reflect.ValueOf(key.IntID())
		case typeOfStr:
			v = reflect.ValueOf(key.StringID())
		default:
			v = reflect.ValueOf(key)
		}
		self.src_val.SetMapIndex(v, d.src_val)
	} else if self.src_kind == reflect.Slice {
		self.src_val.Set(reflect.Append(self.src_val, d.src_val))
	}
}
