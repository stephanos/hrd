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
	var validProps = []ds.Property{
		{Name: "A", Value: "abc"},
		{Name: "B", Value: int64(1)},
	}

	BeforeEach(func() {
		entity = loadEntity{}
		CodecSet.AddMust(entity)
	})

	It("should load an entity from channel properties", func() {
		doc, err := newDocFromInst(&entity)
		Check(err, IsNil)
		Check(doc, NotNil)

		c := newPropertyChannel(validProps)
		err = doc.Load(c)
		Check(err, IsNil)
		Check(c, IsClosed)

		res := (doc.get()).(*loadEntity)
		Check(res.A, Equals, "abc")
		Check(res.B, EqualsNum, 1)
	})

	It("should return an error when loading fails", func() {
		doc, _ := newDocFromInst(&entity)

		c := newPropertyChannel([]ds.Property{{Name: "Invalid"}})
		err := doc.Load(c)

		Check(err, ErrorContains, `cannot load field "Invalid" into a "trafo.loadEntity": no such struct field`)
		Check(c, IsClosed)
	})

	It("should return an error when BeforeLoad fails", func() {
		entity.beforeFunc = func() error {
			return fmt.Errorf("an error")
		}

		doc, _ := newDocFromInst(&entity)
		c := newPropertyChannel(validProps)

		err := doc.Load(c)
		Check(err, ErrorContains, "an error")
		Check(c, IsClosed)

		res := (doc.get()).(*loadEntity)
		Check(res.A, Equals, "")
	})

	It("should return an error when AfterLoad fails", func() {
		entity.afterFunc = func() error {
			return fmt.Errorf("an error")
		}

		doc, _ := newDocFromInst(&entity)
		c := newPropertyChannel(validProps)

		err := doc.Load(c)
		Check(err, ErrorContains, "an error")
		Check(c, IsClosed)

		res := (doc.get()).(*loadEntity)
		Check(res.A, Equals, "abc")
		Check(res.B, EqualsNum, 1)
	})
})

func newPropertyChannel(props []ds.Property) chan ds.Property {
	c := make(chan ds.Property)
	go func() {
		for _, p := range props {
			c <- p
		}
		close(c)
	}()
	return c
}
