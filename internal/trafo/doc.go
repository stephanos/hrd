package trafo

import (
	"fmt"
	"reflect"

	"github.com/101loops/hrd/entity"
	"github.com/101loops/hrd/internal/types"
	"github.com/101loops/structor"
)

// Doc is a reader and writer for a datastore entity.
// It implements the datastore's PropertyLoadSaver.
//
// It is based on:
// https://code.google.com/p/appengine-go/source/browse/appengine/datastore/prop.go
type Doc struct {
	// reference to the entity.
	srcVal reflect.Value

	// codec of the entity.
	codec *structor.Codec
}

func newDoc(srcVal reflect.Value) (*Doc, error) {
	srcType := srcVal.Type()
	srcKind := srcVal.Kind()
	switch srcKind {
	case reflect.Struct:
	case reflect.Ptr:
		srcType = srcVal.Elem().Type()
		srcKind = srcVal.Elem().Kind()
		if srcKind != reflect.Struct {
			return nil, fmt.Errorf("invalid value kind %q (wanted struct pointer)", srcKind)
		}
	default:
		return nil, fmt.Errorf("invalid value kind %q (wanted struct or struct pointer)", srcKind)
	}

	codec, err := getCodec(srcType)
	if err != nil {
		return nil, err
	}

	return &Doc{srcVal, codec}, nil
}

func newDocFromInst(src interface{}) (*Doc, error) {
	return newDoc(reflect.ValueOf(src))
}

func newDocFromType(typ reflect.Type) (*Doc, error) {
	return newDoc(reflect.New(typ.Elem()))
}

// Nil sets the value of the entity to nil.
func (doc *Doc) Nil() {
	dst := doc.val()
	dst.Set(reflect.New(dst.Type()).Elem())
}

// get returns the entity.
func (doc *Doc) get() interface{} {
	return doc.srcVal.Interface()
}

// SetKey assigns a key to the entity.
func (doc *Doc) SetKey(key *types.Key) {
	src := doc.get()

	var parentKey = key.Parent()
	if parentKey != nil {
		id := parentKey.IntID()
		if parent, ok := src.(entity.ParentNumIdentifier); id != 0 && ok {
			parent.SetParent(parentKey.Kind(), id)
		} else {
			sid := parentKey.StringID()
			if parent, ok := src.(entity.ParentTextIdentifier); sid != "" && ok {
				parent.SetParent(parentKey.Kind(), sid)
			}
		}
	}

	id := key.IntID()
	if ident, ok := src.(entity.NumIdentifier); id != 0 && ok {
		ident.SetID(id)
	} else {
		sid := key.StringID()
		if ident, ok := src.(entity.TextIdentifier); sid != "" && ok {
			ident.SetID(sid)
		}
	}
}

func (doc *Doc) val() reflect.Value {
	v := doc.srcVal
	if !v.CanSet() {
		v = doc.srcVal.Elem()
	}
	return v
}
