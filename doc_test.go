package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("Markdown Util", func() {

	It("returns properties", func() {

		doc, err := newDocFromInst(&SimpleModel{
			Num:  42,
			Text: "html",
			Data: []byte("byte"),
		})
		Check(err, IsNil)
		Check(doc, NotNil)

		props, err := doc.toProperties("", []string{}, false)

		Check(err, IsNil)
		Check(props, NotNil)
		Check(props, HasLen, 3)

		Check(*(props[0]), Equals, property{"num", int64(42), false, false})
		Check(*(props[1]), Equals, property{"dat", []byte("byte"), false, false})
		Check(*(props[2]), Equals, property{"html", "html", true, false})

		// # 0

		doc, err = newDocFromInst(&ComplexModel{})
		Check(err, IsNil)
		Check(doc, NotNil)

		props, err = doc.toProperties("", []string{}, false)

		Check(err, IsNil)
		Check(props, NotNil)
		Check(props, HasLen, 1)

		Check(*(props[0]), Equals, property{"tag.val", "", false, false})

		// #1: embedded struct

		doc, err = newDocFromInst(&ComplexModel{
			Pair: Pair{"life", "42"},
		})
		Check(err, IsNil)
		Check(doc, NotNil)

		props, err = doc.toProperties("", []string{}, false)

		Check(err, IsNil)
		Check(props, NotNil)
		Check(props, HasLen, 2)

		Check(*(props[0]), Equals, property{"tag.key", "life", true, false})
		Check(*(props[1]), Equals, property{"tag.val", "42", false, false})

		// #2: embedded struct pointer

		/*
			doc, err = newDoc(&ComplexModel{
				PairPtr: &Pair{"life", "42"},
			})
			Check(err, IsNil)
			Check(doc, IsNil)

			props, err = doc.toProperties("", []string{}, false)

			Check(err, IsNil)
			Check(props, IsNil)
			Check(props, HasLen,   4)

			Check(*(props[0]), Equals, property{"tag.key", "", false, false})
			Check(*(props[1]), Equals, property{"tag.val", "", false, false})
			Check(*(props[2]), Equals, property{"pair.key", "life", true, false})
			Check(*(props[3]), Equals, property{"pair.val", "42", false, false})
		*/

		// #3: embedded slice of structs

		doc, err = newDocFromInst(&ComplexModel{
			Pairs: []Pair{Pair{"Bill", "Bob"}, Pair{"Barb", "Betty"}},
		})
		Check(err, IsNil)
		Check(doc, NotNil)

		props, err = doc.toProperties("", []string{}, false)

		Check(err, IsNil)
		Check(props, NotNil)
		Check(props, HasLen, 5)

		Check(*(props[1]), Equals, property{"tags.key", "Bill", true, true})
		Check(*(props[2]), Equals, property{"tags.val", "Bob", false, true})
		Check(*(props[3]), Equals, property{"tags.key", "Barb", true, true})
		Check(*(props[4]), Equals, property{"tags.val", "Betty", false, true})

		// #4: embedded slice of struct pointers

		/*
			doc, err = newDoc(&ComplexModel{
				PairPtrs: []*Pair{&Pair{"Bill", "Bob"}, &Pair{"Barb", "Betty"}},
			})
			Check(err, IsNil)
			Check(doc, IsNil)

			props, err = doc.toProperties("", []string{}, false)

			Check(err, IsNil)
			Check(props, IsNil)
			Check(props, HasLen,   2)

			Check(*(props[2]), Equals, property{"pairs.key", "Bill", true, true})
			Check(*(props[3]), Equals, property{"pairs.val", "Bob", false, true})
			Check(*(props[4]), Equals, property{"pairs.key", "Barb", true, true})
			Check(*(props[5]), Equals, property{"pairs.val", "Betty", false, true})
		*/
	})
})
