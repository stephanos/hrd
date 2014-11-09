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

	simpleMdls := []*SimpleModel{
		&SimpleModel{id: 1, Text: "text1"}, &SimpleModel{id: 2, Text: "text2"},
		&SimpleModel{id: 3, Text: "text3"}, &SimpleModel{id: 4, Text: "text4"},
	}

	BeforeEach(func() {
		coll = randomColl()

		keys, err := coll.Save().ReqKey().Entities(simpleMdls)
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

	It("deletes multiple entity by ID", func() {
		keys := []*Key{coll.NewNumKey(1), coll.NewNumKey(2)}
		Check(keys[0], ExistsInDatabase)
		Check(keys[1], ExistsInDatabase)

		err := coll.Delete().IDs(keys[0].IntID(), keys[1].IntID())

		Check(err, IsNil)
		Check(keys[0], Not(ExistsInDatabase))
		Check(keys[1], Not(ExistsInDatabase))
	})

	It("deletes all entities", func() {
		_, err := coll.DESTROY()

		Check(err, IsNil)
		// Check(keys, HasLen, 4) TODO
	})
})
