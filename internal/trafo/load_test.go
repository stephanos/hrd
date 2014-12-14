package trafo

import (
	"fmt"
	. "github.com/101loops/bdd"

	ds "appengine/datastore"
)

var _ = Describe("Doc: Load", func() {

	var entity HookEntity
	var validProps = []ds.Property{
		{Name: "A", Value: "abc"},
		{Name: "B", Value: int64(1)},
	}

	BeforeEach(func() {
		entity = HookEntity{}
		CodecSet.AddMust(entity)
	})

	It("should load an entity from channel of properties", func() {
		doc, err := newDocFromInst(&entity)
		Check(err, IsNil)
		Check(doc, NotNil)

		c := newPropertyChannel(validProps)
		err = doc.Load(c)
		Check(err, IsNil)
		Check(c, IsClosed)

		res := (doc.get()).(*HookEntity)
		Check(res.A, Equals, "abc")
		Check(res.B, EqualsNum, 1)
	})

	It("should load an entity with inner struct", func() {
		type InnerModel1 struct {
			Name string
		}
		type InnerModel2 struct {
			Name string `datastore:"name"`
		}
		type MyModel struct {
			InnerModel1
			InnerModel2 `datastore:"inner"`
		}

		CodecSet.AddMust(&MyModel{})
		doc, _ := newDocFromInst(&MyModel{})
		c := newPropertyChannel([]ds.Property{
			{Name: "Name", Value: "him"},
			{Name: "inner.name", Value: "her"},
		})

		err := doc.Load(c)
		Check(err, IsNil)

		res := (doc.get()).(*MyModel)
		Check(res.InnerModel1.Name, Equals, "him")
		Check(res.InnerModel2.Name, Equals, "her")
	})

	// ==== ERRORS

	It("should return an error when loading fails", func() {
		doc, _ := newDocFromInst(&entity)

		c := newPropertyChannel([]ds.Property{{Name: "Invalid"}})
		err := doc.Load(c)

		Check(err, ErrorContains, `cannot load field "Invalid" into a "trafo.HookEntity": no such struct field`)
		Check(c, IsClosed)
	})

	It("should return an error when BeforeLoad fails", func() {
		entity.beforeLoad = func() error {
			return fmt.Errorf("an error")
		}

		doc, _ := newDocFromInst(&entity)
		c := newPropertyChannel(validProps)

		err := doc.Load(c)
		Check(err, ErrorContains, "an error")
		Check(c, IsClosed)

		res := (doc.get()).(*HookEntity)
		Check(res.A, Equals, "")
	})

	It("should return an error when AfterLoad fails", func() {
		entity.afterLoad = func() error {
			return fmt.Errorf("an error")
		}

		doc, _ := newDocFromInst(&entity)
		c := newPropertyChannel(validProps)

		err := doc.Load(c)
		Check(err, ErrorContains, "an error")
		Check(c, IsClosed)

		res := (doc.get()).(*HookEntity)
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
