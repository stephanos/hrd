package hrd

import (
	"appengine/datastore"
	"fmt"
	"github.com/101loops/reflector"
	"reflect"
	"strings"
	"time"
)

// doc is a reader and writer for a datastore entity.
// It implements the datastore's PropertyLoadSaver.
//
// It is based on:
// https://code.google.com/p/appengine-go/source/browse/appengine/datastore/prop.go
type doc struct {
	// reference to entity.
	src_val reflect.Value

	// codec of entity.
	codec *codec
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

func newDoc(src_val reflect.Value) (*doc, error) {
	src_type := src_val.Type()
	src_kind := src_val.Kind()
	switch src_kind {
	case reflect.Struct:
	case reflect.Ptr:
		src_type = src_val.Elem().Type()
		src_kind = src_val.Elem().Kind()
		if src_kind != reflect.Struct {
			return nil, fmt.Errorf("invalid value kind %q (wanted struct pointer)", src_kind)
		}
	default:
		return nil, fmt.Errorf("invalid value kind %q (wanted struct or struct pointer)", src_kind)
	}

	codec, err := getCodec(src_type)
	if err != nil {
		return nil, err
	}

	return &doc{src_val, codec}, nil
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
	return doc.src_val.Interface()
}

// set sets the entity.
func (doc *doc) set(src interface{}) {
	dst := doc.val()
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr && dst.Kind() != reflect.Ptr {
		v = v.Elem()
	}
	dst.Set(v)
}

func (doc *doc) setKey(k *Key) {
	setKey(doc.get(), k)
}

func (doc *doc) val() reflect.Value {
	v := doc.src_val
	if !v.CanSet() {
		v = doc.src_val.Elem()
	}
	return v
}

func (doc *doc) toProperties(prefix string, tags []string, multi bool) (res []*property, err error) {
	var props []*property

	src_val := doc.val()
	for i, t := range doc.codec.byIndex {
		v := src_val.Field(i)
		if !v.IsValid() || !v.CanSet() {
			continue
		}

		name := t.name
		if prefix != "" {
			name = prefix + "." + name
		}

		aggrTags := append(tags, t.tags...)

		// for slice fields (that aren't []byte), save each element
		if v.Kind() == reflect.Slice && v.Type() != typeOfByteSlice {
			for j := 0; j < v.Len(); j++ {
				props, err = itemToProperties(name, aggrTags, true, v.Index(j))
				if err != nil {
					return
				}
				res = append(res, props...)
			}
			continue
		}

		// otherwise, save the field itdoc
		props, err = itemToProperties(name, aggrTags, multi, v)
		if err != nil {
			return
		}
		res = append(res, props...)
	}

	return
}

//
// Note: Save should close the channel when done, even if an error occurred.
func (doc *doc) Save(c chan<- datastore.Property) error {
	defer close(c)

	src := doc.get()

	// event: before save
	if hook, ok := src.(beforeSaver); ok {
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
	if hook, ok := src.(afterSaver); ok {
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
	if hook, ok := dst.(beforeLoader); ok {
		if err := hook.BeforeLoad(); err != nil {
			return err
		}
	}

	if err := datastore.LoadStruct(dst, c); err != nil {
		return err
	}

	// event: after load
	if hook, ok := dst.(afterLoader); ok {
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
			if reflector.IsDefault(v.Interface()) {
				return // ignore complete field if empty
			}
		} else if strings.HasPrefix(tag, "index") {
			indexed = true
			if strings.HasSuffix(tag, ":omitempty") && reflector.IsDefault(v.Interface()) {
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
