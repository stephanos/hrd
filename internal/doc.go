package internal

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/101loops/hrd/entity"
	"github.com/101loops/iszero"
	"github.com/101loops/structor"

	"appengine/datastore"
)

// doc is a reader and writer for a datastore entity.
// It implements the datastore's PropertyLoadSaver.
//
// It is based on:
// https://code.google.com/p/appengine-go/source/browse/appengine/datastore/prop.go
type doc struct {
	// reference to the entity.
	srcVal reflect.Value

	// codec of the entity.
	codec *structor.Codec
}

// property is a name/value pair plus some metadata.
type property struct {
	// name is the property name.
	name string

	// value is the property value.
	value interface{}

	// indexed is whether the datastore indexes this property.
	indexed bool

	// multi is whether the entity can have multiple properties with the same name.
	multi bool
}

func newDoc(srcVal reflect.Value) (*doc, error) {
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

	return &doc{srcVal, codec}, nil
}

func newDocFromInst(src interface{}) (*doc, error) {
	return newDoc(reflect.ValueOf(src))
}

func newDocFromType(typ reflect.Type) (*doc, error) {
	return newDoc(reflect.New(typ.Elem()))
}

// nil sets the value of the entity to nil.
func (doc *doc) nil() {
	dst := doc.val()
	dst.Set(reflect.New(dst.Type()).Elem())
}

// get returns the entity.
func (doc *doc) get() interface{} {
	return doc.srcVal.Interface()
}

func (doc *doc) setKey(k *Key) {
	applyTo(k, doc.get())
}

func applyTo(key *Key, src interface{}) {
	var parentKey = key.Parent()
	if parentKey != nil {
		id := parentKey.IntID()
		if parent, ok := src.(entity.NumParent); id != 0 && ok {
			parent.SetParent(id)
		} else {
			sid := parentKey.StringID()
			if parent, ok := src.(entity.TextParent); sid != "" && ok {
				parent.SetParent(sid)
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

func (doc *doc) val() reflect.Value {
	v := doc.srcVal
	if !v.CanSet() {
		v = doc.srcVal.Elem()
	}
	return v
}

func (doc *doc) toProperties(prefix string, tags []string, multi bool) (res []*property, err error) {
	var props []*property

	srcVal := doc.val()
	for _, fCodec := range doc.codec.Fields() {
		fVal := srcVal.Field(fCodec.Index)
		if !fVal.IsValid() || !fVal.CanSet() {
			continue
		}

		name := fCodec.Label
		if prefix != "" {
			name = prefix + "." + name
		}

		aggrTags := append(tags, fCodec.Tag.Modifiers...)

		// for slice fields (that aren't []byte), save each element
		if fVal.Kind() == reflect.Slice && fVal.Type() != typeOfByteSlice {
			for i := 0; i < fVal.Len(); i++ {
				props, err = itemToProperties(name, aggrTags, true, fVal.Index(i))
				if err != nil {
					return
				}
				res = append(res, props...)
			}
			continue
		}

		// TODO: for map fields, save each element

		props, err = itemToProperties(name, aggrTags, multi, fVal)
		if err != nil {
			return
		}
		res = append(res, props...)
	}

	return
}

// Note: Save should close the channel when done, even if an error occurred.
func (doc *doc) Save(c chan<- datastore.Property) error {
	defer close(c)

	src := doc.get()

	// event: before save
	if hook, ok := src.(entity.BeforeSaver); ok {
		if err := hook.BeforeSave(); err != nil {
			return err
		}
	}

	// export properties
	props, err := doc.toProperties("", []string{""}, false)
	if err != nil {
		return err
	}

	// fill channel
	for _, prop := range props {
		c <- datastore.Property{
			Name:     prop.name,
			Value:    prop.value,
			NoIndex:  !prop.indexed,
			Multiple: prop.multi,
		}
	}

	// event: after save
	if hook, ok := src.(entity.AfterSaver); ok {
		if err := hook.AfterSave(); err != nil {
			close(c)
			return err
		}
	}

	return nil
}

// Note: Load should drain the channel until closed, even if an error occurred.
func (doc *doc) Load(c <-chan datastore.Property) error {

	dst := doc.get()

	// event: before load
	if hook, ok := dst.(entity.BeforeLoader); ok {
		if err := hook.BeforeLoad(); err != nil {
			return err
		}
	}

	if err := datastore.LoadStruct(dst, c); err != nil {
		return err
	}

	// event: after load
	if hook, ok := dst.(entity.AfterLoader); ok {
		if err := hook.AfterLoad(); err != nil {
			return err
		}
	}

	return nil
}

func itemToProperties(name string, tags []string, multi bool, v reflect.Value) (props []*property, err error) {

	// dereference pointers, ignore nil
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	// process tags
	indexed := false
	for _, tag := range tags {
		tag = strings.ToLower(tag)
		if tag == "omitempty" {
			if iszero.Value(v) {
				return // ignore complete field if empty
			}
		} else if strings.HasPrefix(tag, "index") {
			indexed = true
			if strings.HasSuffix(tag, ":omitempty") && iszero.Value(v) {
				indexed = false // ignore index if empty
			}
		} else if tag != "" {
			err = fmt.Errorf("unknown tag %q", tag)
			return
		}
	}

	p := &property{
		name:  name,
		multi: multi,
	}
	props = []*property{p}
	p.indexed = indexed

	// serialize
	switch x := v.Interface().(type) {
	//case *Key:
	//	p.value = x
	case time.Time:
		p.value = x
	case []byte:
		p.indexed = false
		p.value = x
	default:
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			p.value = v.Int()
		case reflect.Bool:
			p.value = v.Bool()
		case reflect.String:
			p.value = v.String()
		case reflect.Float32, reflect.Float64:
			p.value = v.Float()
		case reflect.Struct:
			if !v.CanAddr() {
				return nil, fmt.Errorf("unsupported property %q (unaddressable)", name)
			}
			sub, err := newDocFromInst(v.Addr().Interface())
			if err != nil {
				return nil, fmt.Errorf("unsupported property %q (%v)", name, err)
			}
			return sub.toProperties(name, tags, multi)
		}
	}

	if p.value == nil {
		err = fmt.Errorf("unsupported struct field type %q (unidentifiable)", v.Type())
	}

	return
}
