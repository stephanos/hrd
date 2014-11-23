package trafo

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/101loops/hrd/entity"
	"github.com/101loops/iszero"

	ds "appengine/datastore"
)

// Save saves the entity to datastore properties.
func (d *Doc) Save(c chan<- ds.Property) error {
	defer close(c)

	src := d.get()

	// event hook: before save
	if hook, ok := src.(entity.BeforeSaver); ok {
		if err := hook.BeforeSave(); err != nil {
			return err
		}
	}

	// timestamp
	now := time.Now()
	if ts, ok := src.(entity.CreateTimestamper); ok {
		ts.SetCreatedAt(now)
	}
	if ts, ok := src.(entity.UpdateTimestamper); ok {
		ts.SetUpdatedAt(now)
	}

	// export properties
	props, err := d.toProperties("", []string{""}, false)
	if err != nil {
		return err
	}

	// fill channel
	for _, prop := range props {
		c <- ds.Property{
			Name:     prop.name,
			Value:    prop.value,
			NoIndex:  !prop.indexed,
			Multiple: prop.multi,
		}
	}

	// event hook: after save
	if hook, ok := src.(entity.AfterSaver); ok {
		if err := hook.AfterSave(); err != nil {
			return err
		}
	}

	return nil
}

func (doc *Doc) toProperties(prefix string, tags []string, multi bool) (res []*property, err error) {
	var props []*property

	srcVal := doc.val()
	for _, fCodec := range doc.codec.Fields() {
		fVal := srcVal.Field(fCodec.Index)
		if !fVal.IsValid() || !fVal.CanSet() {
			continue
		}

		name := fCodec.Label
		aggrTags := append(tags, fCodec.Tag.Modifiers...)

		// for slice fields (that aren't []byte), save each element
		if fVal.Kind() == reflect.Slice && fVal.Type() != typeOfByteSlice {
			for i := 0; i < fVal.Len(); i++ {
				props, err = fieldToProps(prefix, name, aggrTags, true, fVal.Index(i))
				if err != nil {
					return
				}
				res = append(res, props...)
			}
			continue
		}

		// TODO: for map fields, save each element

		props, err = fieldToProps(prefix, name, aggrTags, multi, fVal)
		if err != nil {
			return
		}
		res = append(res, props...)
	}

	return
}

func fieldToProps(prefix, name string, tags []string, multi bool, v reflect.Value) (props []*property, err error) {

	// dereference pointers, ignore nil
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	// process tags
	indexed := false
	inlined := false
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
		} else if tag == "inline" {
			inlined = true
		} else if tag != "" {
			err = fmt.Errorf("unknown tag %q", tag)
			return
		}
	}

	p := &property{name: name, multi: multi}
	p.indexed = indexed
	if prefix != "" && !inlined {
		p.name = prefix + propertySeparator + p.name
	}
	props = []*property{p}

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
