package internal

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal/types"
)

var _ = Describe("Put", func() {

	var (
		kind *types.Kind
	)

	BeforeEach(func() {
		kind = randomKind()
		clearCache()
	})

	It("should save an entity without id", func() {
		entity := &MyModel{}
		Check(entity.ID(), EqualsNum, 0)
		Check(entity.UpdatedAt(), IsZero)
		Check(entity.CreatedAt(), IsZero)

		keys, err := Put(kind, entity, false)
		Check(err, IsNil)
		Check(keys, HasLen, 1)

		genID := keys[0].IntID()
		Check(genID, IsGreaterThan, 0)
		Check(entity.ID(), Equals, genID)
		Check(entity.UpdatedAt(), Not(IsZero))
		Check(entity.CreatedAt(), Not(IsZero))
		Check(entity.lifecycle, Equals, []string{"before-save", "after-save"})
	})

	It("should save multiple entities without id", func() {
		entities := []*MyModel{
			&MyModel{}, &MyModel{},
		}

		keys, err := Put(kind, entities, false)
		Check(err, IsNil)
		Check(keys, HasLen, 2)

		Check(keys[0].IntID(), IsGreaterThan, 0)
		Check(entities[0].ID(), EqualsNum, keys[0].IntID())
		Check(keys[1].IntID(), IsGreaterThan, 0)
		Check(entities[1].ID(), EqualsNum, keys[1].IntID())
	})

	It("should save an entity with id", func() {
		entity := &MyModel{}
		entity.SetID(42)

		keys, err := Put(kind, entity, true)
		Check(err, IsNil)
		Check(keys, HasLen, 1)

		Check(entity.ID(), EqualsNum, 42)
		Check(keys[0].IntID(), EqualsNum, 42)
	})

	It("should save multiple entities with id", func() {
		entities := []*MyModel{&MyModel{}, &MyModel{}}
		entities[0].SetID(1)
		entities[1].SetID(2)

		keys, err := Put(kind, entities, true)
		Check(err, IsNil)
		Check(keys, HasLen, 2)

		Check(keys[0].IntID(), EqualsNum, 1)
		Check(keys[1].IntID(), EqualsNum, 2)
	})

	// ==== ERRORS

	It("should not save nil entity", func() {
		keys, err := Put(kind, nil, false)

		Check(keys, IsNil)
		Check(err, NotNil).And(Contains, "must be non-nil")
	})

	// NOTE: other cases of invalid entity/entities are checked inside the trafo package

	It("should not save complete entity without Id", func() {
		entity := &MyModel{}
		keys, err := Put(kind, entity, true)

		Check(keys, IsNil)
		Check(err, NotNil).And(Contains, "is incomplete")
	})

	It("should not save empty entities", func() {
		entities := []*MyModel{}
		keys, err := Put(kind, entities, false)

		Check(keys, IsNil)
		Check(err, NotNil).And(Contains, "no keys provided")
	})
})
