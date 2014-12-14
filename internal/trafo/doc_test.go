package trafo

import (
	"reflect"
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/entity/fixture"
	"github.com/101loops/hrd/internal/types"
)

var _ = Describe("Doc", func() {

	type MyModel struct{}
	type UnknownModel struct{}

	BeforeEach(func() {
		CodecSet.AddMust(MyModel{})
	})

	Context("create from instance", func() {

		It("should create new Doc from known struct", func() {
			doc, err := newDocFromInst(MyModel{})
			Check(err, IsNil)
			Check(doc, NotNil)
		})

		It("should create new Doc from pointer to known struct", func() {
			doc, err := newDocFromInst(&MyModel{})
			Check(err, IsNil)
			Check(doc, NotNil)
		})

		It("should not create new Doc from unknown struct", func() {
			doc, err := newDocFromInst(UnknownModel{})
			Check(doc, IsNil)
			Check(err, ErrorContains, "no registered codec found for type 'trafo.UnknownModel'")
		})

		It("should not create new Doc from pointer to unknown struct", func() {
			doc, err := newDocFromInst(&UnknownModel{})
			Check(doc, IsNil)
			Check(err, ErrorContains, "no registered codec found for type 'trafo.UnknownModel'")
		})

		It("should not create new Doc from non-struct", func() {
			doc, err := newDocFromInst("invalid")
			Check(doc, IsNil)
			Check(err, ErrorContains, `invalid value kind "string" (wanted struct or struct pointer)`)
		})

		It("should not create new Doc from pointer to non-struct", func() {
			invalidEntity := "invalid"
			doc, err := newDocFromInst(&invalidEntity)
			Check(doc, IsNil)
			Check(err, ErrorContains, `invalid value kind "string" (wanted struct pointer)`)
		})
	})

	Context("create from type", func() {

		It("should create new Doc from pointer to known struct", func() {
			doc, err := newDocFromType(reflect.TypeOf(&MyModel{}))
			Check(err, IsNil)
			Check(doc, NotNil)
		})

		It("should not create new Doc from pointer to unknown struct", func() {
			doc, err := newDocFromType(reflect.TypeOf(&UnknownModel{}))
			Check(doc, IsNil)
			Check(err, ErrorContains, "no registered codec found for type 'trafo.UnknownModel'")
		})
	})

	//	It("should set to nil", func() {
	//		entity := &MyModel{}
	//
	//		doc, err := newDocFromInst(&entity)
	//		Check(err, IsNil)
	//		Check(entity, NotNil)
	//
	//		doc.Nil()
	//		Check(entity, IsNil)
	//	})

	Context("set key", func() {

		It("should set key from with numeric id", func() {
			entity := fixture.EntityWithNumID{}
			CodecSet.AddMust(entity)

			doc, _ := newDocFromInst(&entity)
			doc.setKey(types.NewKey("my-kind", "", 42, nil))

			Check(entity.ID(), EqualsNum, 42)
		})

		It("should set key from text id", func() {
			entity := fixture.EntityWithTextID{}
			CodecSet.AddMust(entity)

			doc, _ := newDocFromInst(&entity)
			doc.setKey(types.NewKey("my-kind", "abc", 0, nil))

			Check(entity.ID(), Equals, "abc")
		})

		It("should set key from numeric parent id", func() {
			entity := fixture.EntityWithParentNumID{}
			CodecSet.AddMust(entity)

			doc, _ := newDocFromInst(&entity)
			doc.setKey(types.NewKey("my-kind", "", 1, types.NewKey("my-kind", "", 2, nil)))

			Check(entity.ID(), EqualsNum, 1)
			parentKind, parentID := entity.Parent()
			Check(parentKind, Equals, "my-kind")
			Check(parentID, EqualsNum, 2)
		})

		It("should set key from text parent id", func() {
			entity := fixture.EntityWithParentTextID{}
			CodecSet.AddMust(entity)

			doc, _ := newDocFromInst(&entity)
			doc.setKey(types.NewKey("my-kind", "abc", 0, types.NewKey("my-kind", "xyz", 0, nil)))

			Check(entity.ID(), Equals, "abc")
			parentKind, parentID := entity.Parent()
			Check(parentKind, Equals, "my-kind")
			Check(parentID, Equals, "xyz")
		})
	})
})
