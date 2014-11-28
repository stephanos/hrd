package trafo

import (
	"reflect"
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal/fixture"
	"github.com/101loops/hrd/internal/types"

	ds "appengine/datastore"
)

var _ = Describe("Doc", func() {

	type MyModel struct{}
	type UnknownModel struct{}

	BeforeSuite(func() {
		CodecSet.AddMust(MyModel{})
	})

	Context("create from instance", func() {

		It("known struct", func() {
			doc, err := newDocFromInst(MyModel{})
			Check(err, IsNil)
			Check(doc, NotNil)
		})

		It("pointer to known struct", func() {
			doc, err := newDocFromInst(&MyModel{})
			Check(err, IsNil)
			Check(doc, NotNil)
		})

		It("unknown struct", func() {
			doc, err := newDocFromInst(UnknownModel{})
			Check(doc, IsNil)
			Check(err, NotNil).And(Contains, "no registered codec found for type 'trafo.UnknownModel'")
		})

		It("pointer to unknown struct", func() {
			doc, err := newDocFromInst(&UnknownModel{})
			Check(doc, IsNil)
			Check(err, NotNil).And(Contains, "no registered codec found for type 'trafo.UnknownModel'")
		})
	})

	Context("create from type", func() {

		It("pointer to known struct", func() {
			doc, err := newDocFromType(reflect.TypeOf(&MyModel{}))
			Check(err, IsNil)
			Check(doc, NotNil)
		})

		It("pointer to unknown struct", func() {
			doc, err := newDocFromType(reflect.TypeOf(&UnknownModel{}))
			Check(doc, IsNil)
			Check(err, NotNil).And(Contains, "no registered codec found for type 'trafo.UnknownModel'")
		})
	})

	Context("set key", func() {

		It("with numeric id", func() {
			entity := fixture.EntityWithNumID{}
			CodecSet.AddMust(entity)

			doc, _ := newDocFromInst(&entity)
			doc.setKey(types.NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil)))

			Check(entity.ID(), EqualsNum, 42)
		})

		It("with text id", func() {
			entity := fixture.EntityWithTextID{}
			CodecSet.AddMust(entity)

			doc, _ := newDocFromInst(&entity)
			doc.setKey(types.NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)))

			Check(entity.ID(), Equals, "abc")
		})

		It("with numeric parent id", func() {
			entity := fixture.EntityWithParentNumID{}
			CodecSet.AddMust(entity)

			doc, _ := newDocFromInst(&entity)
			doc.setKey(types.NewKey(ds.NewKey(ctx, "my-kind", "", 1,
				ds.NewKey(ctx, "my-kind", "", 2, nil))))

			Check(entity.ID(), EqualsNum, 1)
			parentKind, parentID := entity.Parent()
			Check(parentKind, Equals, "my-kind")
			Check(parentID, EqualsNum, 2)
		})

		It("with text parent id", func() {
			entity := fixture.EntityWithParentTextID{}
			CodecSet.AddMust(entity)

			doc, _ := newDocFromInst(&entity)
			doc.setKey(types.NewKey(ds.NewKey(ctx, "my-kind", "abc", 0,
				ds.NewKey(ctx, "my-kind", "xyz", 0, nil))))

			Check(entity.ID(), Equals, "abc")
			parentKind, parentID := entity.Parent()
			Check(parentKind, Equals, "my-kind")
			Check(parentID, Equals, "xyz")
		})
	})
})
