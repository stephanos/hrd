package trafo

import (
	"fmt"
	. "github.com/101loops/bdd"

	ds "appengine/datastore"
)

var _ = Describe("Doc: Load", func() {

	var validProps = []ds.Property{
		{Name: "A", Value: "abc"},
		{Name: "B", Value: int64(1)},
	}

	It("should load an entity from properties", func() {
		doc, c, err := load(&HookEntity{}, validProps)
		Check(err, IsNil)
		Check(c, IsClosed)

		res := (doc.get()).(*HookEntity)
		Check(res.A, Equals, "abc")
		Check(res.B, EqualsNum, 1)
	})

	It("should run lifecycle hooks", func() {
		var hooks []string
		entity := &HookEntity{}
		entity.beforeLoad = func() error {
			hooks = append(hooks, "before")
			return nil
		}
		entity.afterLoad = func() error {
			hooks = append(hooks, "after")
			return nil
		}

		_, c, err := load(entity, validProps)
		Check(err, IsNil)
		Check(c, IsClosed)
		Check(hooks, Equals, []string{"before", "after"})
	})

	It("should load an entity with embedded fields from properties", func() {
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

		doc, c, err := load(&MyModel{}, []ds.Property{
			{Name: "Name", Value: "him"},
			{Name: "inner.name", Value: "her"},
		})
		Check(err, IsNil)
		Check(c, IsClosed)

		res := (doc.get()).(*MyModel)
		Check(res.InnerModel1.Name, Equals, "him")
		Check(res.InnerModel2.Name, Equals, "her")
	})

	// ==== ERRORS

	It("should return an error when loading fails", func() {
		entity := HookEntity{}
		_, c, err := load(&entity, []ds.Property{{Name: "Invalid"}})
		Check(err, ErrorContains, `cannot load field "Invalid" into a "trafo.HookEntity": no such struct field`)
		Check(c, IsClosed)
	})

	It("should return an error when BeforeLoad fails", func() {
		entity := HookEntity{}
		entity.beforeLoad = func() error {
			return fmt.Errorf("an error")
		}

		doc, c, err := load(&entity, validProps)
		Check(err, ErrorContains, "an error")
		Check(c, IsClosed)

		res := (doc.get()).(*HookEntity)
		Check(res.A, Equals, "")
	})

	It("should return an error when AfterLoad fails", func() {
		entity := HookEntity{}
		entity.afterLoad = func() error {
			return fmt.Errorf("an error")
		}

		doc, c, err := load(&entity, validProps)
		Check(err, ErrorContains, "an error")
		Check(c, IsClosed)

		res := (doc.get()).(*HookEntity)
		Check(res.A, Equals, "abc")
		Check(res.B, EqualsNum, 1)
	})
})

func load(src interface{}, props []ds.Property) (*Doc, chan ds.Property, error) {
	CodecSet.AddMust(src)
	doc, err := newDocFromInst(src)
	if err != nil {
		return nil, nil, err
	}

	c := make(chan ds.Property)
	go func() {
		for _, p := range props {
			c <- p
		}
		close(c)
	}()

	err = doc.Load(c)
	return doc, c, err
}
