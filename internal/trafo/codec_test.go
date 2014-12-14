package trafo

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/structor"
)

var _ = Describe("Codec", func() {

	It("should return codec", func() {
		type MyModel struct {
			Num   int
			Name  string `datastore:"Label"`
			Label bool   `datastore:"-"`
		}

		entity := &MyModel{}
		err := CodecSet.Add(entity)
		Check(err, IsNil)

		var codec *structor.Codec
		codec, err = getCodec(entity)
		Check(err, IsNil)
		Check(codec, NotNil)
		Check(codec.Complete, IsTrue)
	})

	// ==== ERRORS

	It("should return error for invalid codec", func() {
		codec, err := getCodec("invalid-type")

		Check(codec, IsNil)
		Check(err, ErrorContains, `value is not a struct, struct pointer or reflect.Type - but "string"`)
	})

	It("should reject field beginning with invalid character", func() {
		type InvalidModel struct {
			InvalidName string `datastore:"$invalid-name"`
		}
		err := CodecSet.Add(InvalidModel{})
		Check(err, ErrorContains, `field "InvalidName" begins with invalid character '$'`)

		// as sub-field:
		type Wrapper struct {
			Inner InvalidModel
		}
		err = CodecSet.Add(Wrapper{})
		Check(err, ErrorContains, `field "InvalidName" begins with invalid character '$'`)
	})

	It("should reject fields containing invalid character", func() {
		type InvalidModel struct {
			InvalidName string `datastore:"invalid@name"`
		}
		err := CodecSet.Add(InvalidModel{})
		Check(err, ErrorContains, `field "InvalidName" contains invalid character '@'`)

		// as sub-field:
		type Wrapper struct {
			Inner InvalidModel
		}
		err = CodecSet.Add(Wrapper{})
		Check(err, ErrorContains, `field "InvalidName" contains invalid character '@'`)
	})

	It("should reject duplicate field names", func() {
		type InvalidModel struct {
			ID1 string `datastore:"id"`
			ID2 string `datastore:"id"`
		}
		err := CodecSet.Add(InvalidModel{})
		Check(err, ErrorContains, `duplicate field name "id"`)

		// from sub-field:
		type InnerModel struct {
			ID string `datastore:"id"`
		}
		type MyModel struct {
			ID string `datastore:"id"`
			InnerModel
		}
		err = CodecSet.Add(MyModel{})
		Check(err, ErrorContains, `duplicate field name "id"`)
	})

	It("should reject invalid field type", func() {
		type InvalidModel struct {
			Ptr *string
		}
		err := CodecSet.Add(InvalidModel{})
		Check(err, ErrorContains, `field "Ptr" has invalid type 'pointer'`)
	})

	It("should reject invalid map key type", func() {
		type InvalidModel struct {
			Map map[int]string
		}
		err := CodecSet.Add(InvalidModel{})
		Check(err, ErrorContains, `field "Map" has invalid map key type 'int' - only 'string' is allowed`)
	})

	It("should reject recursive struct", func() {
		type InvalidModel struct {
			Recursive []InvalidModel
		}

		err := CodecSet.Add(InvalidModel{})
		Check(err, ErrorContains, `recursive struct at field "Recursive"`)
	})

	It("should reject slice of slices", func() {
		type Model struct {
			InnerSlice []string
		}
		type InvalidModel struct {
			OuterSlide []Model
		}

		err := CodecSet.Add(InvalidModel{})
		Check(err, ErrorContains, `field "OuterSlide" leads to a slice of slices`)
	})
})
