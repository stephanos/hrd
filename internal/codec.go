package internal

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/101loops/structor"
)

var (
	// CodecSet is a set of all registered entity codecs.
	CodecSet *structor.Set

	typeOfStr       = reflect.TypeOf("")
	typeOfByteSlice = reflect.TypeOf([]byte(nil))
	typeOfTime      = reflect.TypeOf(time.Time{})
)

func init() {
	newCodecSet()
}

func newCodecSet() {
	CodecSet = structor.NewSet("datastore")
	CodecSet.SetValidateFunc(validateCodec)
}

// GetCodec returns an entity's codec.
// The entity must be been added to the codec set beforehand.
func GetCodec(entity interface{}) (*structor.Codec, error) {
	codec, err := CodecSet.Get(entity)
	if err != nil {
		return nil, err
	}

	return codec, nil
}

func validateCodec(codecSet *structor.Set, codec *structor.Codec) error {
	labels := make(map[string]bool, 0)

	for _, field := range codec.Fields() {
		fType := field.Type

		// field ignored?
		if strings.HasPrefix(field.Label, "_") {
			continue
		}

		// valid field name?
		if err := validateFieldName(field.Label); err != nil {
			return fmt.Errorf("field %q has invalid name (%v)", field.Name, err)
		}

		label := strings.ToLower(field.Label)
		if _, ok := labels[label]; ok {
			return fmt.Errorf("duplicate field name %q", label)
		}
		labels[label] = true

		// valid field type?
		if err := validateFieldType(field.Type); err != nil {
			return fmt.Errorf("field %q has invalid type (%v)", field.Name, err)
		}

		if field.KeyType != nil {
			keyType := *field.KeyType
			if keyType != typeOfStr {
				return fmt.Errorf("field %q has invalid map key type '%v' - only 'string' is allowed", field.Name, keyType)
			}
		}

		// valid sub-codec?
		var innerType *reflect.Type
		if fType.Kind() == reflect.Struct {
			if fType != typeOfTime {
				innerType = &fType
			}
		} else if field.ElemType == nil {
			if fType != typeOfByteSlice {
				innerType = field.ElemType
			}
		}

		if innerType != nil {
			subCodec, err := codecSet.Get(*innerType)
			if err != nil {
				return fmt.Errorf("error processing field %q (%v)", field.Name, err)
			}

			if !subCodec.Complete() {
				return fmt.Errorf("recursive struct at field %q", field.Name)
			}

			hasSlice := false
			for _, subField := range subCodec.Fields() {
				label := strings.ToLower(label + "." + subField.Label)
				if _, ok := labels[label]; ok {
					return fmt.Errorf("duplicate field name %q", label)
				}
				labels[label] = true

				if subField.Type.Kind() == reflect.Slice {
					hasSlice = true
				}
			}

			if fType.Kind() == reflect.Slice && hasSlice {
				return fmt.Errorf("field %q leads to a slice of slices", field.Name)
			}
		}
	}

	return nil
}

func validateFieldType(typ reflect.Type) error {
	if typ.Kind() == reflect.Ptr {
		return fmt.Errorf("field is a pointer")
	}

	return nil
}

func validateFieldName(name string) error {
	if name == "" {
		return fmt.Errorf("field name is empty")
	}

	if strings.Contains(name, ".") {
		return fmt.Errorf("field name contains '.'")
	}

	first := true
	for _, char := range name {
		if first {
			first = false
			if char != '_' && !unicode.IsLetter(char) {
				return fmt.Errorf("field name begins with invalid character %q", char)
			}
		} else {
			if char != '_' && !unicode.IsLetter(char) && !unicode.IsDigit(char) {
				return fmt.Errorf("field name contains invalid character %q", char)
			}
		}
	}

	return nil
}
