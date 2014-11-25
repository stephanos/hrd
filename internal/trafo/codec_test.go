package trafo

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/structor"
)

var _ = Describe("Codec", func() {

	It("return simple codec", func() {
		entity := &SimpleModel{}
		err := CodecSet.Add(entity)
		Check(err, IsNil)

		var codec *structor.Codec
		codec, err = getCodec(entity)
		Check(err, IsNil)
		Check(codec, NotNil)
		Check(codec.Complete(), IsTrue)

		fieldNames := codec.FieldNames()
		Check(fieldNames, HasLen, 7)
		Check(fieldNames, Equals, []string{"NumID", "CreatedTime", "UpdatedTime", "Num", "Data", "Text", "Time"})

		fields := codec.Fields()
		Check(fields, HasLen, 7)
	})

	It("return complex codec", func() {
		entity := &ComplexModel{}

		err := CodecSet.Add(entity)
		Check(err, IsNil)

		var codec *structor.Codec
		codec, err = getCodec(entity)
		Check(err, IsNil)
		Check(codec, NotNil)
	})

	// ==== ERRORS

	It("return error for invalid codec", func() {
		codec, err := getCodec("invalid-type")

		Check(codec, IsNil)
		Check(err, NotNil).And(Contains, `value is not a struct, struct pointer or reflect.Type - but "string"`)
	})

	It("rejects field beginning with invalid character", func() {
		type InvalidModel struct {
			InvalidName string `datastore:"$invalid-name"`
		}
		err := CodecSet.Add(InvalidModel{})
		Check(err, NotNil).And(Contains, `field "InvalidName" begins with invalid character '$'`)

		// as sub-field:
		type Wrapper struct {
			Inner InvalidModel
		}
		err = CodecSet.Add(Wrapper{})
		Check(err, NotNil).And(Contains, `field "InvalidName" begins with invalid character '$'`)
	})

	It("rejects fields containing invalid character", func() {
		type InvalidModel struct {
			InvalidName string `datastore:"invalid@name"`
		}
		err := CodecSet.Add(InvalidModel{})
		Check(err, NotNil).And(Contains, `field "InvalidName" contains invalid character '@'`)

		// as sub-field:
		type Wrapper struct {
			Inner InvalidModel
		}
		err = CodecSet.Add(Wrapper{})
		Check(err, NotNil).And(Contains, `field "InvalidName" contains invalid character '@'`)
	})

	It("rejects duplicate field names", func() {
		type InvalidModel struct {
			ID1 string `datastore:"id"`
			ID2 string `datastore:"id"`
		}
		err := CodecSet.Add(InvalidModel{})
		Check(err, NotNil).And(Contains, `duplicate field name "id"`)
	})

	It("rejects invalid field type", func() {
		type InvalidModel struct {
			Ptr *string
		}
		err := CodecSet.Add(InvalidModel{})
		Check(err, NotNil).And(Contains, `field "Ptr" has invalid type 'pointer'`)
	})

	It("rejects invalid map key type", func() {
		type InvalidModel struct {
			Map map[int]string
		}
		err := CodecSet.Add(InvalidModel{})
		Check(err, NotNil).And(Contains, `field "Map" has invalid map key type 'int' - only 'string' is allowed`)
	})

	It("rejects recursive struct", func() {
		type InvalidModel struct {
			Recursive []InvalidModel
		}

		err := CodecSet.Add(InvalidModel{})
		Check(err, NotNil).And(Contains, `recursive struct at field "Recursive"`)
	})

	//	It("rejects slice of slices", func() {
	//		type Model struct {
	//			InnerSlice []string
	//		}
	//		type InvalidModel struct {
	//			OuterSlide []Model
	//		}
	//
	//		err := CodecSet.Add(InvalidModel{})
	//		Check(err, NotNil).And(Contains, `TODO`)
	//	})
})
