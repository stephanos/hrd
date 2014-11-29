package trafo

import (
	. "github.com/101loops/bdd"

	ds "appengine/datastore"
)

var _ = Describe("Doc: Load", func() {

	It("should load an entity from channel properties", func() {
		type MyModel struct {
			A string
			B int
		}
		CodecSet.AddMust(MyModel{})

		doc, err := newDocFromInst(&MyModel{})
		Check(err, IsNil)
		Check(doc, NotNil)

		c := make(chan ds.Property)
		go func() {
			c <- ds.Property{Name: "A", Value: "abc"}
			c <- ds.Property{Name: "B", Value: int64(1)}
			close(c)
		}()

		err = doc.Load(c)
		Check(err, IsNil)
		Check(c, IsClosed)

		res := (doc.get()).(*MyModel)
		Check(res.A, Equals, "abc")
		Check(res.B, EqualsNum, 1)
	})
})
