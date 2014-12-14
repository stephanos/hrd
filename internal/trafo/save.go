package trafo

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/101loops/hrd/entity"
	"github.com/101loops/hrd/internal/types"
	"github.com/101loops/iszero"

	ae "appengine"
	ds "appengine/datastore"
)

var (
	nowFunc = func() time.Time {
		return time.Now()
	}
)

// Save saves the entity to datastore properties.
func (doc *Doc) Save(ctx ae.Context) (props []*ds.Property, err error) {
	src := doc.get()

	// event hook: before save
	if hook, ok := src.(entity.BeforeSaver); ok {
		if err = hook.BeforeSave(); err != nil {
			return
		}
	}

	// timestamp
	now := nowFunc()
	if ts, ok := src.(entity.CreateTimestamper); ok {
		ts.SetCreatedAt(now)
	}
	if ts, ok := src.(entity.UpdateTimestamper); ok {
		ts.SetUpdatedAt(now)
	}

	// export properties
	props, err = doc.toProperties(ctx, "", []string{""}, false)
	if err != nil {
		return
	}

	// event hook: after save
	if hook, ok := src.(entity.AfterSaver); ok {
		if err = hook.AfterSave(); err != nil {
			return
		}
	}

	return
}

func (doc *Doc) toProperties(ctx ae.Context, prefix string, tags []string, multi bool) (res []*ds.Property, err error) {
	var props []*ds.Property
	res = make([]*ds.Property, 0)

	srcVal := doc.val()
	for _, fCodec := range doc.codec.Fields {
		fVal := srcVal.Field(fCodec.Index)
		if !fVal.IsValid() || !fVal.CanSet() {
			continue
		}

		name := fCodec.Attrs["label"].(string)
		if name == "-" {
			continue
		}
		aggrTags := append(tags, fCodec.Tag.Modifiers()...)

		// for slice fields (that aren't []byte), save each element
		if fVal.Kind() == reflect.Slice && fVal.Type() != typeOfByteSlice {
			for i := 0; i < fVal.Len(); i++ {
				props, err = fieldToProps(ctx, prefix, name, aggrTags, true, fVal.Index(i))
				if err != nil {
					return
				}
				res = append(res, props...)
			}
			continue
		}

		props, err = fieldToProps(ctx, prefix, name, aggrTags, multi, fVal)
		if err != nil {
			return
		}
		res = append(res, props...)
	}

	return
}

func fieldToProps(ctx ae.Context, prefix, name string, tags []string, multi bool, v reflect.Value) (props []*ds.Property, err error) {

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

	p := &ds.Property{Name: name, NoIndex: !indexed, Multiple: multi}
	if prefix != "" {
		p.Name = prefix + propertySeparator + p.Name
	}
	props = []*ds.Property{p}

	// serialize
	switch x := v.Interface().(type) {
	case *ds.Key:
		p.Value = x
	case types.DSKeyConverter:
		p.Value = x.ToDSKey(ctx)
	case time.Time:
		p.Value = x
	case ae.BlobKey:
		p.Value = x
	case ae.GeoPoint:
		p.Value = x
	case []byte:
		p.Value = x
		p.NoIndex = true
	default:
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			p.Value = v.Int()
		case reflect.Bool:
			p.Value = v.Bool()
		case reflect.String:
			p.Value = v.String()
		case reflect.Float32, reflect.Float64:
			p.Value = v.Float()
		case reflect.Struct:
			if !v.CanAddr() {
				return nil, fmt.Errorf("unsupported property %q (unaddressable)", name)
			}
			sub, err := newDocFromInst(v.Addr().Interface())
			if err != nil {
				return nil, fmt.Errorf("unsupported property %q (%v)", name, err)
			}
			return sub.toProperties(ctx, name, tags, multi)
		}
	}

	if p.Value == nil {
		err = fmt.Errorf("unsupported struct field type %q (unidentifiable)", v.Type())
	}

	return
}
