package hrd

import (
	"fmt"
	"github.com/101loops/reflector"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"
)

// codec describes how to convert a struct to and from a sequence of properties.
type codec struct {
	// byIndex gives the tagCodec for the i'th field.
	byIndex map[int]tagCodec
	// byName gives the field codec for the tagCodec with the passed name.
	byName map[string]fieldCodec
	// hasSlice is whether a struct or any of its nested or embedded structs
	// has a slice-typed field (other than []byte).
	hasSlice bool
	// complete is whether the codec is complete.
	// An incomplete codec may be encountered when walking a recursive struct.
	complete bool
}

type tagCodec struct {
	name string
	tags []string
}

type fieldCodec struct {
	index    int
	subcodec *codec
}

var (
	typeOfByteSlice = reflect.TypeOf([]byte(nil))
	typeOfTime      = reflect.TypeOf(time.Time{})

	codecDictMutex sync.Mutex
	codecDict      = make(map[reflect.Type]*codec)
	codecTags      = []string{"datastore"}
)

func getCodec(obj interface{}) (*codec, error) {
	codecDictMutex.Lock()
	defer codecDictMutex.Unlock()
	return getCodecLocked(obj)
}

// Note: codecDictMutex must be held when calling this function.
func getCodecLocked(obj interface{}) (*codec, error) {
	refl, err := reflector.NewStructCodec(obj)
	if err != nil {
		return nil, err
	}
	return getCodecStructLocked(refl)
}

// Note: codecDictMutex must be held when calling this function.
func getCodecStructLocked(refl *reflector.StructCodec) (*codec, error) {

	t := refl.Type()

	// lookup and return cached codec, if any
	ret, ok := codecDict[t]
	if ok {
		return ret, nil
	}

	defer func() {
		if ret != nil {
			delete(codecDict, t)
		}
	}()

	fields, err := refl.FieldCodecs(codecTags)
	if err != nil {
		return nil, err
	}

	// create new codec
	ret = &codec{
		byIndex: make(map[int]tagCodec),
		byName:  make(map[string]fieldCodec),
	}
	// Add c to the codecDict map before we are sure it is good.
	// If t is a recursive type, it needs to find the incomplete entry for itself in the map.
	codecDict[t] = ret

	for _, f := range fields {

		if strings.HasPrefix(f.Label, "_") {
			continue // skip fields starting with '_' (e.g. '_id')
		}

		// validate field label
		nameErrMsg := validatePropertyName(f.Label)
		if nameErrMsg != "" {
			return nil, fmt.Errorf("field %q has invalid name (%v)", f.Name, nameErrMsg)
		}

		// validate field type
		typeErrMsg := validatePropertyType(f.Type)
		if typeErrMsg != "" {
			return nil, fmt.Errorf("field %q has invalid type (%v)", f.Name, typeErrMsg)
		}

		// is field a sub-struct?
		subType, fIsSlice := reflect.Type(nil), false
		switch f.Type.Kind() {
		//case reflect.Ptr:
		//	subType = f.Type.Elem()
		case reflect.Struct:
			subType = f.Type
		case reflect.Slice:
			sliceElem := f.Type.Elem()
			sliceElemKind := sliceElem.Kind()
			if sliceElemKind == reflect.Ptr {
				subType = sliceElem.Elem()
			} else if sliceElemKind == reflect.Struct {
				subType = sliceElem
			}
			fIsSlice = f.Type != typeOfByteSlice
			ret.hasSlice = ret.hasSlice || fIsSlice
		}

		if subType != nil && subType != typeOfTime {

			// process sub-struct
			sub, err := getCodecLocked(subType)
			if err != nil {
				return nil, fmt.Errorf("error processing field %q (%v)", f.Name, err)
			}
			if !sub.complete {
				return nil, fmt.Errorf("recursive struct at field %q", f.Name)
			}
			if fIsSlice && sub.hasSlice {
				return nil, fmt.Errorf("field %q leads to a slice of slices", f.Name)
			}
			ret.hasSlice = ret.hasSlice || sub.hasSlice
			for relName := range sub.byName {
				absName := f.Label + "." + relName
				if _, ok := ret.byName[absName]; ok {
					return nil, fmt.Errorf("duplicate property name %q", absName)
				}
				ret.byName[absName] = fieldCodec{index: f.Index, subcodec: sub}
			}
		} else {

			// process non-sub-struct
			if _, ok := ret.byName[f.Label]; ok {
				return nil, fmt.Errorf("duplicate property name %q", f.Label)
			}
			ret.byName[f.Label] = fieldCodec{index: f.Index}
		}

		ret.byIndex[f.Index] = tagCodec{
			name: f.Label,
			tags: f.Tags,
		}
	}
	ret.complete = true

	return ret, nil
}

func validatePropertyType(typ reflect.Type) string {
	return "" // TODO ?
}

func validatePropertyName(name string) string {
	if name == "" {
		return "property name is empty"
	}

	if strings.Contains(name, ".") {
		return "property name contains '.'"
	}

	first := true
	for _, c := range name {
		if first {
			first = false
			if c != '_' && !unicode.IsLetter(c) {
				return fmt.Sprintf("property name begins with invalid character %q", c)
			}
		} else {
			if c != '_' && !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return fmt.Sprintf("property name contains invalid character %q", c)
			}
		}
	}

	return ""
}
