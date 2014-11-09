package internal

import (
	"time"
	. "github.com/101loops/bdd"
	"github.com/101loops/structor"
)

var _ = Describe("Codec", func() {

	It("evaluate simple codec", func() {
		type SimpleModel struct {
			id     int64
			Ignore string    `datastore:"-"`
			Num    int64     `datastore:"num"`
			Data   []byte    `datastore:",index"`
			Text   string    `datastore:"html,index"`
			Time   time.Time `datastore:"timing,index,omitempty"`
		}

		entity := &SimpleModel{}

		err := CodecSet.Add(entity)
		Check(err, IsNil)

		var code *structor.Codec
		code, err = GetCodec(entity)
		Check(err, IsNil)
		Check(code, NotNil)
		Check(code.Complete(), IsTrue)

		fieldNames := code.FieldNames()
		Check(fieldNames, HasLen, 4)
		Check(fieldNames, Equals, []string{"Num", "Data", "Text", "Time"})

		fields := code.Fields()
		Check(fields, HasLen, 4)
	})

	It("evaluate complex codec", func() {
		type Pair struct {
			Key string `datastore:"key,index,omitempty"`
			Val string
		}

		type ComplexModel struct {
			Pair Pair `datastore:"tag"`
			//PairPtr  *Pair   `datastore:"pair"`
			Pairs []Pair `datastore:"tags"`
			//PairPtrs []*Pair `datastore:"pairs"`
		}

		entity := &ComplexModel{}

		err := CodecSet.Add(entity)
		Check(err, IsNil)

		var code *structor.Codec
		code, err = GetCodec(entity)
		Check(err, IsNil)
		Check(code, NotNil)
	})

	// ==== ERRORS

	It("rejects invalid field names", func() {
		type InvalidCodec1 struct {
			InvalidName string `datastore:"$invalid-name"`
		}
		err := CodecSet.Add(InvalidCodec1{})
		Check(err, NotNil).And(Contains, `field "InvalidName" has invalid name (begins with invalid character '$')`)

		type InvalidCodec2 struct {
			InvalidName string `datastore:"invalid@name"`
		}
		err = CodecSet.Add(InvalidCodec2{})
		Check(err, NotNil).And(Contains, `field "InvalidName" has invalid name (contains invalid character '@')`)
	})

	It("rejects duplicate field names", func() {
		type InvalidCodec struct {
			ID1 string `datastore:"id"`
			ID2 string `datastore:"id"`
		}
		err := CodecSet.Add(InvalidCodec{})
		Check(err, NotNil).And(Contains, `duplicate field name "id"`)
	})

	It("rejects invalid field type", func() {
		type InvalidCodec struct {
			Ptr *string
		}
		err := CodecSet.Add(InvalidCodec{})
		Check(err, NotNil).And(Contains, `field "Ptr" has invalid type (field is a pointer)`)
	})

	It("rejects invalid map key type", func() {
		type InvalidCodec struct {
			Map map[int]string
		}
		err := CodecSet.Add(InvalidCodec{})
		Check(err, NotNil).And(Contains, `field "Map" has invalid map key type 'int' - only 'string' is allowed`)
	})

})
