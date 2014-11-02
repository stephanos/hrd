package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("HRD Load/Save/Delete", func() {

	With("default settings", func() {
		basicTests()
	})

	With("w/o global cache", func() {
		basicTests(NoGlobalCache)
	})
})

func basicTests(opts ...Opt) {

	const (
		setID = 42
	)

	var (
		genID  int64
		coll   *Collection
		loader *Loader
		entity *SimpleModel
	)

	BeforeEach(func() {
		if coll == nil {
			coll = randomColl()
		}

		loader = coll.Load(opts...)
		clearCache()
	})

	It("saves an entity without id", func() {
		entity = &SimpleModel{Text: "text", Data: []byte{1, 2, 3}}
		key, err := coll.Save().Entity(entity)

		Check(err, IsNil)
		Check(key, NotNil)

		genID = key.IntID()
		Check(genID, IsGreaterThan, 0)
		Check(entity.ID(), Equals, genID)
		Check(entity.lifecycle, Equals, []string{"before-save", "after-save"})
		Check(entity.updatedAt, Not(IsZero))
		Check(entity.createdAt, Not(IsZero))
	})

	It("saves an entity with id", func() {
		entity.SetID(setID)
		key, err := coll.Save(CompleteKeys).Entity(entity)

		Check(err, IsNil)
		Check(key, NotNil)

		Check(key.IntID(), EqualsNum, setID)
		Check(entity.ID(), EqualsNum, setID)
	})

	It("loads an entity", func() {
		var entity *SimpleModel
		key, err := loader.ID(setID).GetOne(&entity)

		Check(err, IsNil)
		Check(key, NotNil)

		Check(entity, NotNil)
		Check(entity.Text, Equals, "text")
		Check(entity.ID(), EqualsNum, setID)
		// Check(entity.lifecycle, Equals, []string{"before-load", "after-load"}) TODO
	})

	It("loads all entities into slice of struct pointers", func() {
		var entities []*SimpleModel
		keys, err := coll.Load().IDs(1, setID, genID).GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 3)
		Check(entities, HasLen, 3)

		Check(keys[0].IntID(), EqualsNum, 1)
		Check(keys[0].Exists(), IsFalse)

		Check(keys[1].IntID(), EqualsNum, setID)
		Check(entities[1].Text, Equals, "text")

		Check(keys[2].IntID(), EqualsNum, genID)
		Check(entities[2].Text, Equals, "text")
	})

	It("loads all entities into map of Key to struct pointers", func() {
		var entities map[*Key]*SimpleModel
		keys, err := coll.Load().IDs(1, setID, genID).GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 3)
		Check(entities, HasLen, 3)
	})

	It("loads all entities into map of int64 to struct pointers", func() {
		var entities map[int64]*SimpleModel
		keys, err := coll.Load().IDs(1, setID, genID).GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 3)
		Check(entities, HasLen, 3)
	})

	It("deletes an entity", func() {
		err := coll.Delete().ID(setID)
		Check(err, IsNil)

		var entity *SimpleModel
		key, err := loader.ID(setID).GetOne(&entity)
		Check(err, IsNil)
		Check(key, IsNil)
		Check(entity, IsNil)
	})
}
