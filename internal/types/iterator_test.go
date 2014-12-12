package types

import (
	. "github.com/101loops/bdd"

	ae "appengine"
	ds "appengine/datastore"
)

var _ = Describe("Iterator", func() {

	var (
		qry *Query
	)

	BeforeEach(func() {
		qry = NewQuery("my-kind")
	})

	It("should create a new iterator", func() {
		it := NewIterator(ctx, qry)
		Check(it.inner, NotNil)
	})

	It("should return next entity", func() {
		_, err := NewIterator(ctx, qry).Next(func(_ ae.Context) ds.PropertyLoadSaver {
			return nil
		})
		Check(err, Equals, ds.Done)
	})

	Context("cursor", func() {

		It("should return a cursor", func() {
			cursor, err := NewIterator(ctx, qry).Cursor()

			Check(err, IsNil)
			Check(cursor, IsZero)
		})

		It("should return an error for invalid query", func() {
			qry.Filter = []Filter{Filter{"age !=", 18}}
			cursor, err := NewIterator(ctx, qry).Cursor()

			Check(cursor, IsZero)
			Check(err, HasOccurred).And(Contains, "invalid operator")
		})
	})
})
