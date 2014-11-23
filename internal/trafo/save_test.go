package trafo

import (
	"time"
	. "github.com/101loops/bdd"
	ds "appengine/datastore"
)

var _ = Describe("Doc", func() {

	Context("returns properties", func() {

		It("simple model", func() {
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

		It("complex model", func() {
			doc, err := newDocFromInst(&ComplexModel{})
			Check(err, IsNil)
			Check(doc, NotNil)

			props, err := doc.toProperties(ctx, "", []string{}, false)

			Check(err, IsNil)
			Check(props, NotNil)
			Check(props, HasLen, 1)

			Check(*props[0], Equals, ds.Property{"tag.Val", "", true, false})
		})

		It("complex model with inner struct", func() {
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

		It("complex model with slice of structs", func() {
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
})
