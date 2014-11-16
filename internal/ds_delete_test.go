package internal

import (
	. "github.com/101loops/bdd"

	ds "appengine/datastore"
)

var _ = Describe("DSDelete", func() {

	var (
		kind Kind
	)

	var entities []interface{}

	BeforeEach(func() {
		kind = randomKind()

		entities = []interface{}{
			&SimpleModel{id: 1}, &SimpleModel{id: 2},
			&SimpleModel{id: 3}, &SimpleModel{id: 4},
		}
		keys, err := DSPut(kind, entities, true)
		Check(err, IsNil)
		Check(keys, HasLen, 4)

		clearCache()
	})

	It("delete entities by key", func() {
		key := ds.NewKey(ctx, kind.Name(), "", 1, nil)
		Check(existsInDB(key), IsTrue)

		err := DSDeleteKeys(kind, newKeys(key))

		Check(err, IsNil)
		Check(existsInDB(key), IsFalse)
	})

	It("deletes multiple entities by key", func() {
		keys := []*ds.Key{
			ds.NewKey(ctx, kind.Name(), "", 1, nil),
			ds.NewKey(ctx, kind.Name(), "", 2, nil),
		}
		Check(existsInDB(keys[0]), IsTrue)
		Check(existsInDB(keys[1]), IsTrue)

		err := DSDeleteKeys(kind, newKeys(keys...))

		Check(err, IsNil)
		Check(existsInDB(keys[0]), IsFalse)
		Check(existsInDB(keys[1]), IsFalse)
	})

	It("deletes entity", func() {
		key := ds.NewKey(ctx, kind.Name(), "", 1, nil)
		Check(existsInDB(key), IsTrue)

		err := DSDelete(kind, entities[0], false)

		Check(err, IsNil)
		Check(existsInDB(key), IsFalse)
	})

	It("deletes slice of entities", func() {
		keys := []*ds.Key{
			ds.NewKey(ctx, kind.Name(), "", 1, nil),
			ds.NewKey(ctx, kind.Name(), "", 2, nil),
		}
		Check(existsInDB(keys[0]), IsTrue)
		Check(existsInDB(keys[1]), IsTrue)

		err := DSDelete(kind, entities[0:2], true)

		Check(err, IsNil)
		Check(existsInDB(keys[0]), IsFalse)
		Check(existsInDB(keys[1]), IsFalse)
	})

	It("deletes map of entities", func() {
		keys := []*ds.Key{
			ds.NewKey(ctx, kind.Name(), "", 1, nil),
			ds.NewKey(ctx, kind.Name(), "", 2, nil),
		}
		Check(existsInDB(keys[0]), IsTrue)
		Check(existsInDB(keys[1]), IsTrue)

		entityMap := map[string]interface{}{"a": entities[0], "b": entities[1]}
		err := DSDelete(kind, entityMap, true)

		Check(err, IsNil)
		Check(existsInDB(keys[0]), IsFalse)
		Check(existsInDB(keys[1]), IsFalse)
	})

	// ==== ERRORS

	It("does not delete invalid entity", func() {
		var entity string
		err := DSDelete(kind, entity, false)

		Check(err, NotNil).And(Contains, `value type "string" does not provide ID()`)
	})

	It("does not delete invalid entities", func() {
		entities := []string{"a", "b", "c"}
		err := DSDelete(kind, entities, true)

		Check(err, NotNil).And(Contains, `value type "string" does not provide ID()`)
	})
})
