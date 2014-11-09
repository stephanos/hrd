package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("HRD Delete", func() {

	const (
		setID = 42
	)

	var (
		coll *Collection
	)

	var entities []interface{}

	BeforeEach(func() {
		coll = randomColl()

		entities = []interface{}{
			&SimpleModel{id: 1, Text: "text1"},
			&SimpleModel{id: 2, Text: "text2"},
			&ChildModel{id: "a", parentID: 1, parentKind: coll.name},
			&ChildModel{id: "b", parentID: 2, parentKind: coll.name},
		}
		keys, err := coll.Save(CompleteKeys).Entities(entities)
		Check(err, IsNil)
		Check(keys, HasLen, 4)

		clearCache()
	})

	It("deletes an entity by key", func() {
		key := coll.NewNumKey(1)
		Check(key, ExistsInDatabase)

		err := coll.Delete().Key(key)

		Check(err, IsNil)
		Check(key, Not(ExistsInDatabase))
	})

	It("deletes multiple entities by key", func() {
		keys := []*Key{coll.NewNumKey(1), coll.NewNumKey(2)}
		Check(keys[0], ExistsInDatabase)
		Check(keys[1], ExistsInDatabase)

		err := coll.Delete().Keys(keys)

		Check(err, IsNil)
		Check(keys[0], Not(ExistsInDatabase))
		Check(keys[1], Not(ExistsInDatabase))
	})

	It("deletes an entity by ID", func() {
		key := coll.NewNumKey(1)
		Check(key, ExistsInDatabase)

		err := coll.Delete().ID(key.IntID())

		Check(err, IsNil)
		Check(key, Not(ExistsInDatabase))
	})

	It("deletes an entity by text ID", func() {
		parentKey := coll.NewNumKey(1)
		key := coll.NewTextKey("a", parentKey)
		Check(key, ExistsInDatabase)

		err := coll.Delete().TextID(key.StringID(), parentKey)

		Check(err, IsNil)
		Check(key, Not(ExistsInDatabase))
	})

	It("deletes multiple entity by ID", func() {
		keys := []*Key{coll.NewNumKey(1), coll.NewNumKey(2)}
		Check(keys[0], ExistsInDatabase)
		Check(keys[1], ExistsInDatabase)

		err := coll.Delete().IDs(keys[0].IntID(), keys[1].IntID())

		Check(err, IsNil)
		Check(keys[0], Not(ExistsInDatabase))
		Check(keys[1], Not(ExistsInDatabase))
	})

	It("deletes entity", func() {
		key := coll.NewNumKey(1)
		Check(key, ExistsInDatabase)

		err := coll.Delete().Entity(entities[0])

		Check(err, IsNil)
		Check(key, Not(ExistsInDatabase))
	})

	It("deletes slice of entities", func() {
		keys := []*Key{coll.NewNumKey(1), coll.NewNumKey(2)}
		Check(keys[0], ExistsInDatabase)
		Check(keys[1], ExistsInDatabase)

		err := coll.Delete().Entities(entities[0:2])

		Check(err, IsNil)
		Check(keys[0], Not(ExistsInDatabase))
		Check(keys[1], Not(ExistsInDatabase))
	})

	It("deletes map of entities", func() {
		keys := []*Key{coll.NewNumKey(1), coll.NewNumKey(2)}
		Check(keys[0], ExistsInDatabase)
		Check(keys[1], ExistsInDatabase)

		entityMap := map[string]interface{}{"a": entities[0], "b": entities[1]}
		err := coll.Delete().Entities(entityMap)

		Check(err, IsNil)
		Check(keys[0], Not(ExistsInDatabase))
		Check(keys[1], Not(ExistsInDatabase))
	})

	It("deletes all entities", func() {
		_, err := coll.DESTROY()

		Check(err, IsNil)
		// Check(keys, HasLen, 4) TODO
	})

	// ==== ERRORS

	It("does not delete invalid entity", func() {
		var entity string
		err := coll.Delete().Entity(entity)

		Check(err, NotNil).And(Contains, `value type "string" does not provide ID()`)
	})

	It("does not delete invalid entities", func() {
		entities := []string{"a", "b", "c"}
		err := coll.Delete().Entities(entities)

		Check(err, NotNil).And(Contains, `value type "string" does not provide ID()`)
	})
})
