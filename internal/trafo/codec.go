package trafo

import (
	"errors"
	"fmt"
	"reflect"
	"time"
	"unicode"

	"github.com/101loops/structor"

	ae "appengine"
	ds "appengine/datastore"
)

const (
	propertySeparator = "."
)

var (
	// CodecSet is a set of all registered entity codecs.
	CodecSet *structor.Set

	errFieldIgnored = errors.New("field ignored")

	typeOfStr       = reflect.TypeOf("")
	typeOfDSKey     = reflect.TypeOf((*ds.Key)(nil))
	typeOfByteSlice = reflect.TypeOf([]byte(nil))
	typeOfTime      = reflect.TypeOf(time.Time{})
	typeOfGeoPoint  = reflect.TypeOf(ae.GeoPoint{})
)

func init() {
	newCodecSet()
}

func newCodecSet() {
	CodecSet = structor.NewSet("datastore")
	CodecSet.SetValidateFunc(validateCodec)
}

// getCodec returns an entity's codec.
// The entity must be been added to the codec set beforehand.
func getCodec(entity interface{}) (*structor.Codec, error) {
	codec, err := CodecSet.Get(entity)
	if err != nil {
		return nil, err
	}

	return codec, nil
}

func validateCodec(_ *structor.Set, codec *structor.Codec) error {
	labels := make(map[string]bool, 0)

	for _, field := range codec.Fields {

		err := validateFieldName(labels, field)
		if err == errFieldIgnored {
			continue
		}
		if err != nil {
			return fmt.Errorf("field %q %v", field.Name, err)
		}

		err = validateFieldType(field)
		if err != nil {
			return fmt.Errorf("field %q %v", field.Name, err)
		}

		if err := validateSubField(labels, field); err != nil {
			return err
		}
	}

	return nil
}

func validateFieldType(field *structor.FieldCodec) error {
	if field.Type.Kind() == reflect.Ptr && field.Type != typeOfDSKey {
		return fmt.Errorf("has invalid type 'pointer'")
	}

	if field.KeyType != nil {
		keyType := *field.KeyType
		if keyType != typeOfStr {
			return fmt.Errorf("has invalid map key type '%v' - only 'string' is allowed", keyType)
		}
	}

	return nil
}

func validateFieldName(labels map[string]bool, field *structor.FieldCodec) error {
	label := calcLabel(field)

	if label == "" {
		return nil
	}
	if label == "-" {
		return errFieldIgnored
	}

	first := true
	for _, char := range label {
		if first {
			first = false
			if char != '_' && !unicode.IsLetter(char) {
				return fmt.Errorf("begins with invalid character %q", char)
			}
		} else {
			if char != '_' && !unicode.IsLetter(char) && !unicode.IsDigit(char) {
				return fmt.Errorf("contains invalid character %q", char)
			}
		}
	}

	if _, ok := labels[label]; ok {
		return fmt.Errorf("has duplicate field name %q", label)
	}
	labels[label] = true

	return nil
}

func validateSubField(labels map[string]bool, parentField *structor.FieldCodec) error {
	subType := subTypeOf(parentField.Type, parentField.ElemType)
	if subType == nil || subType.Kind() != reflect.Struct {
		return nil
	}

	subCodec, _ := CodecSet.Get(subType)
	if !subCodec.Complete {
		return fmt.Errorf("recursive struct at field %q", parentField.Name)
	}

	parentLabel := parentField.Attrs["label"].(string)
	if parentLabel != "" {
		parentLabel = parentLabel + propertySeparator
	}

	hasSlice := false
	for _, subField := range subCodec.Fields {
		if subField.Type.Kind() == reflect.Slice {
			hasSlice = true
		}

		subLabel := parentLabel + calcLabel(subField)
		if _, ok := labels[subLabel]; ok {
			return fmt.Errorf("duplicate field name %q", subLabel)
		}
		labels[subLabel] = true
	}

	if parentField.Type.Kind() == reflect.Slice && hasSlice {
		return fmt.Errorf("field %q leads to a slice of slices", parentField.Name)
	}

	return nil
}

func subTypeOf(fieldType reflect.Type, elemType *reflect.Type) reflect.Type {
	if fieldType.Kind() == reflect.Struct {
		if fieldType != typeOfTime {
			return fieldType
		}
	} else if elemType != nil {
		if fieldType != typeOfByteSlice && fieldType != typeOfGeoPoint {
			return *elemType
		}
	}

	return nil
}

func calcLabel(f *structor.FieldCodec) string {
	label := f.Tag.Values[0]
	if label == "" {
		if !f.Anonymous {
			label = f.Name
		}
	}
	f.Attrs["label"] = label
	return label
}
