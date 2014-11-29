package trafo

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/entity/fixture"
	"github.com/101loops/hrd/internal/types"

	ds "appengine/datastore"
)

var _ = Describe("DocList", func() {

	var (
		kind     *types.Kind
		keys     []*types.Key
		entities []*fixture.EntityWithNumID
	)

	type UnknownModel struct{}
	type InvalidModel struct{}

	BeforeEach(func() {
		CodecSet.AddMust(&InvalidModel{})
		kind = types.NewKind(ctx, "my-kind")

		keys = make([]*types.Key, 4)
		entities = make([]*fixture.EntityWithNumID, 4)
		for i := int64(0); i < 4; i++ {
			entity := &fixture.EntityWithNumID{}
			entity.SetID(i + 1)
			entities[i] = entity
			keys[i] = types.NewKey(ds.NewKey(ctx, kind.Name, "", entity.ID(), nil))
		}
	})

	Context("readable list", func() {

		It("should create list from struct pointer", func() {
			list, err := NewReadableDocList(kind, entities[0])
			Check(err, IsNil)
			Check(list, NotNil)
			Check(list.list, HasLen, 1)
		})

		It("should create list from slice of struct pointers", func() {
			list, err := NewReadableDocList(kind, entities[0:2])
			Check(err, IsNil)
			Check(list, NotNil)
			Check(list.list, HasLen, 2)
		})

		It("should not create list not from nil value", func() {
			list, err := NewReadableDocList(kind, nil)
			Check(list, IsNil)
			Check(err, NotNil).And(Contains, "value must be non-nil")
		})

		It("should not create list from unknown struct pointer", func() {
			list, err := NewReadableDocList(kind, &UnknownModel{})
			Check(list, IsNil)
			Check(err, NotNil).And(Contains, "no registered codec found for type 'trafo.UnknownModel'")
		})

		It("should not create list from invalid entity pointer", func() {
			list, err := NewReadableDocList(kind, &InvalidModel{})
			Check(list, IsNil)
			Check(err, NotNil).And(Contains, `value type "*trafo.InvalidModel" does not provide ID()`)
		})
	})

	Context("writeable list", func() {

		It("should create list from struct pointer", func() {
			list, err := NewWriteableDocList(&(entities[0]), keys[0:1], false)
			Check(err, IsNil)
			Check(list, NotNil)
			Check(list.list, HasLen, 1)
		})

		It("should create list from slice", func() {
			var entitySlice []*fixture.EntityWithNumID
			list, err := NewWriteableDocList(&entitySlice, keys, true)

			Check(err, IsNil)
			Check(list, NotNil)
			Check(list.list, HasLen, 4)
		})

		It("should create list from map", func() {
			var entityMap map[string]*fixture.EntityWithNumID
			list, err := NewWriteableDocList(&entityMap, keys, true)

			Check(err, IsNil)
			Check(list, NotNil)
			Check(list.list, HasLen, 4)
		})

		It("should create list from map", func() {
			var entityMap map[string]*fixture.EntityWithNumID
			list, err := NewWriteableDocList(&entityMap, keys, true)

			Check(err, IsNil)
			Check(list, NotNil)
			Check(list.list, HasLen, 4)
		})

		It("should not create list from non-pointer", func() {
			list, err := NewWriteableDocList(*(entities[0]), keys[0:1], false)
			Check(list, IsNil)
			Check(err, NotNil).And(Contains, `invalid value kind "struct" (wanted non-nil pointer)`)
		})

		It("should not create list from nil pointer", func() {
			var entity *fixture.EntityWithNumID
			list, err := NewWriteableDocList(entity, keys[0:1], false)
			Check(list, IsNil)
			Check(err, NotNil).And(Contains, `invalid value kind "ptr" (wanted non-nil pointer)`)
		})
	})
})
