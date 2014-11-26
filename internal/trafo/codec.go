package trafo

import (
	"fmt"
	"reflect"
	"strings"
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

	for _, field := range codec.Fields() {

		label, err := validateFieldName(labels, field)
		if err != nil {
			return fmt.Errorf("field %q %v", field.Name, err)
		}

		fType, err := validateFieldType(field)
		if err != nil {
			return fmt.Errorf("field %q %v", field.Name, err)
		}

		subType := subTypeOf(fType, field.ElemType)
		if subType != nil && subType.Kind() == reflect.Struct {
			if err := validateSubType(labels, label, field.Name, subType, fType.Kind()); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateFieldType(field *structor.FieldCodec) (reflect.Type, error) {
	fType := field.Type

	if fType.Kind() == reflect.Ptr && fType != typeOfDSKey {
		return nil, fmt.Errorf("has invalid type 'pointer'")
	}

	if field.KeyType != nil {
		keyType := *field.KeyType
		if keyType != typeOfStr {
			return nil, fmt.Errorf("has invalid map key type '%v' - only 'string' is allowed", keyType)
		}
	}

	return fType, nil
}

func validateFieldName(labels map[string]bool, field *structor.FieldCodec) (string, error) {
	label := strings.ToLower(field.Label)

	first := true
	for _, char := range label {
		if first {
			first = false
			if char != '_' && !unicode.IsLetter(char) {
				return "", fmt.Errorf("begins with invalid character %q", char)
			}
		} else {
			if char != '_' && !unicode.IsLetter(char) && !unicode.IsDigit(char) {
				return "", fmt.Errorf("contains invalid character %q", char)
			}
		}
	}

	if _, ok := labels[label]; ok {
		return "", fmt.Errorf("duplicate field name %q", label)
	}
	labels[label] = true

	return label, nil
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

func validateSubType(labels map[string]bool, fLabel string, fName string,
	subType reflect.Type, parentKind reflect.Kind) error {

	subCodec, err := CodecSet.Get(subType)
	if err != nil {
		return fmt.Errorf("error processing field %q (%v)", fName, err)
	}

	if !subCodec.Complete() {
		return fmt.Errorf("recursive struct at field %q", fName)
	}

	hasSlice := false
	for _, subField := range subCodec.Fields() {
		subLabel := strings.ToLower(subField.Label)
		if !subField.Tag.HasModifier("inline") {
			subLabel = fLabel + propertySeparator + subLabel
		}
		if _, ok := labels[subLabel]; ok {
			return fmt.Errorf("duplicate field name %q", subLabel)
		}
		labels[subLabel] = true

		if subField.Type.Kind() == reflect.Slice {
			hasSlice = true
		}
	}

	if parentKind == reflect.Slice && hasSlice {
		return fmt.Errorf("field %q leads to a slice of slices", fName)
	}

	return nil
}
