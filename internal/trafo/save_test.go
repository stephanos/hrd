package trafo

import (
	"time"
	. "github.com/101loops/bdd"

	ae "appengine"
	ds "appengine/datastore"
)

var _ = Describe("Doc Save", func() {

	toProps := func(src interface{}) ([]*ds.Property, error) {
		CodecSet.AddMust(src)
		doc, err := newDocFromInst(src)
		if err != nil {
			panic(err)
		}
		return doc.Save(ctx)
	}

	It("saves primitives", func() {
		type MyModel struct {
			I   int
			I8  int8
			I16 int16
			I32 int32
			I64 int64
			B   bool
			S   string
			F32 float32
			F64 float64
		}

		props, err := toProps(&MyModel{
			int(1), int8(2), int16(3), int32(4), int64(5), true, "test", float32(1.0), float64(2.0),
		})
		Check(err, IsNil)
		Check(props, NotNil).And(HasLen, 9)

		Check(*props[0], Equals, ds.Property{"I", int64(1), true, false})
		Check(*props[1], Equals, ds.Property{"I8", int64(2), true, false})
		Check(*props[2], Equals, ds.Property{"I16", int64(3), true, false})
		Check(*props[3], Equals, ds.Property{"I32", int64(4), true, false})
		Check(*props[4], Equals, ds.Property{"I64", int64(5), true, false})
		Check(*props[5], Equals, ds.Property{"B", true, true, false})
		Check(*props[6], Equals, ds.Property{"S", "test", true, false})
		Check(*props[7], Equals, ds.Property{"F32", float64(1.0), true, false})
		Check(*props[8], Equals, ds.Property{"F64", float64(2.0), true, false})
	})

	It("saves known complex types", func() {
		type MyModel struct {
			B  []byte
			T  time.Time
			BK ae.BlobKey
			GP ae.GeoPoint
		}

		entity := &MyModel{
			[]byte("test"), time.Now(), ae.BlobKey("bkey"), ae.GeoPoint{1, 2},
		}
		props, err := toProps(entity)
		Check(err, IsNil)
		Check(props, NotNil).And(HasLen, 4)

		Check(*props[0], Equals, ds.Property{"B", entity.B, true, false})
		Check(*props[1], Equals, ds.Property{"T", entity.T, true, false})
		Check(*props[2], Equals, ds.Property{"BK", entity.BK, true, false})
		Check(*props[3], Equals, ds.Property{"GP", entity.GP, true, false})
	})

	It("saves simple model", func() {
		doc, err := newDocFromInst(&SimpleModel{
			Num:  42,
			Text: "html",
			Data: []byte("byte"),
		})
		Check(err, IsNil)
		Check(doc, NotNil)

		props, err := doc.toProperties(ctx, "", []string{}, false)

		Check(err, IsNil)
		Check(props, NotNil)
		Check(props, HasLen, 6)

		Check(*props[0], Equals, ds.Property{"id", int64(0), true, false})
		Check(*props[1], Equals, ds.Property{"created_at", time.Time{}, false, false})
		Check(*props[2], Equals, ds.Property{"updated_at", time.Time{}, false, false})
		Check(*props[3], Equals, ds.Property{"num", int64(42), true, false})
		Check(*props[4], Equals, ds.Property{"Data", []byte("byte"), true, false})
		Check(*props[5], Equals, ds.Property{"html", "html", false, false})
	})

	It("saves complex model", func() {
		doc, err := newDocFromInst(&ComplexModel{})
		Check(err, IsNil)
		Check(doc, NotNil)

		props, err := doc.toProperties(ctx, "", []string{}, false)

		Check(err, IsNil)
		Check(props, NotNil)
		Check(props, HasLen, 1)

		Check(*props[0], Equals, ds.Property{"tag.Val", "", true, false})
	})

	It("saves complex model with inner struct", func() {
		doc, err := newDocFromInst(&ComplexModel{
			Pair: Pair{"life", "42"},
		})
		Check(err, IsNil)
		Check(doc, NotNil)

		props, err := doc.toProperties(ctx, "", []string{}, false)

		Check(err, IsNil)
		Check(props, NotNil)
		Check(props, HasLen, 2)

		Check(*props[0], Equals, ds.Property{"tag.key", "life", false, false})
		Check(*props[1], Equals, ds.Property{"tag.Val", "42", true, false})
	})

	/*
		It("complex model with inner struct pointer", func() {

			doc, err = newDoc(&ComplexModel{
				PairPtr: &Pair{"life", "42"},
			})
			Check(err, IsNil)
			Check(doc, IsNil)

			props, err = doc.toProperties("", []string{}, false)

			Check(err, IsNil)
			Check(props, IsNil)
			Check(props, HasLen,   4)

			Check(*(props[0]), Equals, ds.Property{"tag.key", "", false, false})
			Check(*(props[1]), Equals, ds.Property{"tag.val", "", false, false})
			Check(*(props[2]), Equals, ds.Property{"pair.key", "life", true, false})
			Check(*(props[3]), Equals, ds.Property{"pair.val", "42", false, false})
		})
	*/

	It("saves complex model with slice of structs", func() {
		doc, err := newDocFromInst(&ComplexModel{
			Pairs: []Pair{Pair{"Bill", "Bob"}, Pair{"Barb", "Betty"}},
		})
		Check(err, IsNil)
		Check(doc, NotNil)

		props, err := doc.toProperties(ctx, "", []string{}, false)

		Check(err, IsNil)
		Check(props, NotNil)
		Check(props, HasLen, 5)

		Check(*props[1], Equals, ds.Property{"tags.key", "Bill", false, true})
		Check(*props[2], Equals, ds.Property{"tags.Val", "Bob", true, true})
		Check(*props[3], Equals, ds.Property{"tags.key", "Barb", false, true})
		Check(*props[4], Equals, ds.Property{"tags.Val", "Betty", true, true})
	})

	/*
		It("complex model with slice of struct pointers", func() {
			doc, err = newDoc(&ComplexModel{
				PairPtrs: []*Pair{&Pair{"Bill", "Bob"}, &Pair{"Barb", "Betty"}},
			})
			Check(err, IsNil)
			Check(doc, IsNil)

			props, err = doc.toProperties("", []string{}, false)

			Check(err, IsNil)
			Check(props, IsNil)
			Check(props, HasLen,   2)

			Check(*(props[2]), Equals, ds.Property{"pairs.key", "Bill", true, true})
			Check(*(props[3]), Equals, ds.Property{"pairs.val", "Bob", false, true})
			Check(*(props[4]), Equals, ds.Property{"pairs.key", "Barb", true, true})
			Check(*(props[5]), Equals, ds.Property{"pairs.val", "Betty", false, true})
		})
	*/
})
