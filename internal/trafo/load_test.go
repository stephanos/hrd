package trafo

import (
	"fmt"
	. "github.com/101loops/bdd"

	ds "appengine/datastore"
)

type loadEntity struct {
	A string
	B int

	beforeFunc func() error
	afterFunc  func() error
}

func (l *loadEntity) BeforeLoad() error {
	if l.beforeFunc != nil {
		return l.beforeFunc()
	}
	return nil
}

func (l *loadEntity) AfterLoad() error {
	if l.afterFunc != nil {
		return l.afterFunc()
	}
	return nil
}

var _ = Describe("Doc: Load", func() {

	var entity loadEntity

	BeforeEach(func() {
		entity = loadEntity{}
		CodecSet.AddMust(entity)
	})

	It("should load an entity from channel properties", func() {
		doc, err := newDocFromInst(&entity)
		Check(err, IsNil)
		Check(doc, NotNil)

		c := newPropertyChannel()
		err = doc.Load(c)
		Check(err, IsNil)
		Check(c, IsClosed)

		res := (doc.get()).(*loadEntity)
		Check(res.A, Equals, "abc")
		Check(res.B, EqualsNum, 1)
	})

	It("should return an error when BeforeLoad errs", func() {
		entity.beforeFunc = func() error {
			return fmt.Errorf("an error")
		}

		doc, err := newDocFromInst(&entity)
		Check(err, IsNil)
		Check(doc, NotNil)

		c := newPropertyChannel()
		err = doc.Load(c)
		Check(err, ErrorContains, "an error")
		Check(c, IsClosed)

		res := (doc.get()).(*loadEntity)
		Check(res.A, Equals, "")
	})

	It("should return an error when AfterLoad errs", func() {
		entity.afterFunc = func() error {
			return fmt.Errorf("an error")
		}

		doc, err := newDocFromInst(&entity)
		Check(err, IsNil)
		Check(doc, NotNil)

		c := newPropertyChannel()
		err = doc.Load(c)
		Check(err, ErrorContains, "an error")
		Check(c, IsClosed)

		res := (doc.get()).(*loadEntity)
		Check(res.A, Equals, "abc")
		Check(res.B, EqualsNum, 1)
	})
})

func newPropertyChannel() chan ds.Property {
	c := make(chan ds.Property)
	go func() {
		c <- ds.Property{Name: "A", Value: "abc"}
		c <- ds.Property{Name: "B", Value: int64(1)}
		close(c)
	}()
	return c
}
