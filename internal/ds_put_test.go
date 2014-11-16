package internal

import . "github.com/101loops/bdd"

var _ = Describe("DSPut", func() {

	var (
		kind Kind
	)

	BeforeEach(func() {
		kind = randomKind()
		clearCache()
	})

	It("saves an entity without id", func() {
		entity := &SimpleModel{}

		keys, err := DSPut(kind, entity, false)
		Check(err, IsNil)
		Check(keys, HasLen, 1)

		genID := keys[0].IntID()
		Check(genID, IsGreaterThan, 0)
		Check(entity.ID(), Equals, genID)
		Check(entity.lifecycle, Equals, []string{"before-save", "after-save"})
	})

	It("saves multiple entities without id", func() {
		entities := []*SimpleModel{
			&SimpleModel{}, &SimpleModel{},
		}

		keys, err := DSPut(kind, entities, false)
		Check(err, IsNil)
		Check(keys, HasLen, 2)

		Check(keys[0].IntID(), IsGreaterThan, 0)
		Check(entities[0].ID(), EqualsNum, keys[0].IntID())
		Check(keys[1].IntID(), IsGreaterThan, 0)
		Check(entities[1].ID(), EqualsNum, keys[1].IntID())
	})

	It("saves an entity with id", func() {
		entity := &SimpleModel{id: 42}

		keys, err := DSPut(kind, entity, true)
		Check(err, IsNil)
		Check(keys, HasLen, 1)

		Check(entity.ID(), EqualsNum, 42)
		Check(keys[0].IntID(), EqualsNum, 42)
	})

	It("saves multiple entities with id", func() {
		entities := []*SimpleModel{
			&SimpleModel{id: 1}, &SimpleModel{id: 2},
		}

		keys, err := DSPut(kind, entities, true)
		Check(err, IsNil)
		Check(keys, HasLen, 2)

		Check(keys[0].IntID(), EqualsNum, 1)
		Check(keys[1].IntID(), EqualsNum, 2)
	})

	It("saves multiple entities with parents", func() {
		entities := []*ChildModel{
			&ChildModel{id: "a", parentID: 42, parentKind: kind.Name()},
			&ChildModel{id: "b", parentID: 42, parentKind: kind.Name()},
		}
		keys, err := DSPut(kind, entities, true)

		Check(err, IsNil)
		Check(keys, HasLen, 2)
		Check(keys[0].StringID(), Equals, "a")
		Check(keys[1].StringID(), Equals, "b")
	})

	// ==== ERRORS

	It("does not save nil entity", func() {
		keys, err := DSPut(kind, nil, false)

		Check(keys, IsNil)
		Check(err, NotNil).And(Contains, "must be non-nil")
	})

	It("does not save non-struct entity", func() {
		keys, err := DSPut(kind, 42, false)

		Check(keys, IsNil)
		Check(err, NotNil).And(Contains, "invalid value kind").And(Contains, "int")
	})

	It("does not save entity without ID() and 42()", func() {
		invalidMdl := &InvalidModel{}
		keys, err := DSPut(kind, invalidMdl, false)

		Check(keys, IsNil)
		Check(err, NotNil).And(Contains, "does not provide ID")
	})

	It("does not save complete entity without Id", func() {
		entity := &SimpleModel{}
		keys, err := DSPut(kind, entity, true)

		Check(keys, IsNil)
		Check(err, NotNil).And(Contains, "is incomplete")
	})

	It("does not save empty entities", func() {
		entities := []*SimpleModel{}
		keys, err := DSPut(kind, entities, false)

		Check(keys, IsNil)
		Check(err, NotNil).And(Contains, "no keys provided")
	})
})
