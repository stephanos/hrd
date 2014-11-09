package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("HRD Load", func() {

	With("default settings", func() {
		loadTests()
	})

	With("w/o global cache", func() {
		loadTests(NoGlobalCache)
	})
})

func loadTests(opts ...Opt) {

	var (
		coll   *Collection
		loader *Loader
	)

	simpleMdls := []*SimpleModel{
		&SimpleModel{id: 1, Text: "text1"}, &SimpleModel{id: 2, Text: "text2"},
		&SimpleModel{id: 3, Text: "text3"}, &SimpleModel{id: 4, Text: "text4"},
	}

	BeforeEach(func() {
		if coll == nil {
			coll = randomColl()
		}

		keys, err := coll.Save().ReqKey().Entities(simpleMdls)
		Check(err, IsNil)
		Check(keys, HasLen, 4)
		Check(keys[0].IntID(), EqualsNum, 1)
		Check(keys[1].IntID(), EqualsNum, 2)
		Check(keys[2].IntID(), EqualsNum, 3)
		Check(keys[3].IntID(), EqualsNum, 4)

		clearCache()
	})

	BeforeEach(func() {
		if coll == nil {
			coll = randomColl()
		}

		loader = coll.Load(opts...)
		clearCache()
	})

	It("loads an entity", func() {
		var entity *SimpleModel
		validate := func(key *Key, err error) {
			Check(err, IsNil)
			Check(key, NotNil)

			Check(entity, NotNil)
			Check(entity.ID(), EqualsNum, 1)
			Check(entity.Text, Equals, "text1")
			// Check(entity.lifecycle, Equals, []string{"before-load", "after-load"}) TODO
		}

		key, err := loader.ID(1).GetOne(&entity)
		validate(key, err)

		key, err = loader.Key(key).GetOne(&entity)
		validate(key, err)
	})

	It("loads all entities into slice of struct pointers", func() {
		var entities []*SimpleModel
		keys, err := coll.Load().IDs(666, 2, 3).GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 3)
		Check(entities, HasLen, 3)

		Check(keys[0].IntID(), EqualsNum, 666)
		Check(keys[0].Exists(), IsFalse)

		Check(keys[1].IntID(), EqualsNum, 2)
		Check(entities[1].Text, Equals, "text2")

		Check(keys[2].IntID(), EqualsNum, 3)
		Check(entities[2].Text, Equals, "text3")
	})

	It("loads all entities into map of Key to struct pointers", func() {
		var entities map[*Key]*SimpleModel
		keys, err := coll.Load().IDs(1, 2, 3).GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 3)
		Check(entities, HasLen, 3)
	})

	It("loads all entities into map of int64 to struct pointers", func() {
		var entities map[int64]*SimpleModel
		keys, err := coll.Load().IDs(1, 2, 3).GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 3)
		Check(entities, HasLen, 3)
	})

	// ==== ERRORS

	It("does not load entity into invalid type", func() {
		var entity string
		key, err := coll.Load().ID(1).GetOne(entity)

		Check(key, IsNil)
		Check(err, NotNil).And("invalid value kind").And(Contains, "string")
	})

	It("does not load entity into non-pointer struct", func() {
		var entity SimpleModel
		key, err := coll.Load().ID(1).GetOne(entity)

		Check(key, IsNil)
		Check(err, NotNil).And(Contains, "invalid value kind").And(Contains, "struct")
	})

	It("does not load entity into non-reference struct", func() {
		var entity *SimpleModel
		key, err := coll.Load().ID(1).GetOne(entity)

		Check(key, IsNil)
		Check(err, NotNil).And(Contains, "invalid value kind").And(Contains, "ptr")
	})

	It("does not load entities into map with invalid key", func() {
		var entities map[bool]*SimpleModel
		keys, err := coll.Load().IDs(1, 2, 3).GetAll(&entities)

		Check(keys, IsEmpty)
		Check(err, NotNil).And(Contains, "invalid value key")
	})

	It("does not accept key for different collection", func() {
		invalidKey := store.NewNumKey("wrong-kind", 1)

		var entity *SimpleModel
		key, err := coll.Load().Key(invalidKey).GetOne(&entity)

		Check(key, IsNil)
		Check(entity, IsNil)
		Check(err, NotNil).And(Contains, "invalid key kind 'wrong-kind'")

		var keys []*Key
		var entities []*SimpleModel
		keys, err = coll.Load().Keys(invalidKey).GetAll(&entities)

		Check(keys, IsNil)
		Check(entities, IsNil)
		Check(err, NotNil).And(Contains, "invalid key kind 'wrong-kind'")
	})

	It("does not load non-existing entity", func() {
		var entity *SimpleModel
		key, err := coll.Load().ID(666).GetOne(&entity)

		Check(err, IsNil)
		Check(key, IsNil)
		Check(entity, IsNil)
	})
}
