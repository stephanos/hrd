package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("Cache", func() {

	var (
		coll *Collection
	)

	BeforeEach(func() {
		coll = randomColl()
		clearCache()

		entities := []*SimpleModel{
			&SimpleModel{id: 1, Text: "text1"},
			&SimpleModel{id: 2, Text: "text2"},
		}

		keys, err := coll.Save().ReqKey().Entities(entities)
		Check(err, IsNil)
		Check(keys, HasLen, 2)
	})

	Describe("CRUD", func() {
		crudTest := func(source string, flags ...Flag) {
			var entities []*SimpleModel
			keys, err := coll.Load().Flags(flags...).IDs(1, 2).GetAll(&entities)

			Check(err, IsNil)
			Check(keys, HasLen, 2)
			Check(keys[0].source, Equals, source)
			Check(keys[1].source, Equals, source)

			Check(entities, HasLen, 2)
			Check(entities[0].Text, Equals, "text1")
			Check(entities[1].Text, Equals, "text2")
		}

		It("uses memory cache", func() {
			crudTest(SOURCE_MEMORY)
		})

		It("uses memcache", func() {
			crudTest(SOURCE_MEMCACHE, NO_LOCAL_CACHE)
		})

		It("ignores cache", func() {
			crudTest(SOURCE_DATASTORE, NO_CACHE)
		})
	})

	//	Describe("Query", func() {
	//		It("uses hybrid query", func() {
	//			var entities []*SimpleModel
	//			keys, _, err := coll.Query().GetAll(&entities)
	//
	////			var entities []*SimpleModel
	////			keys, err := coll.Load().IDs(1, 2).GetAll(&entities)
	//
	//			Check(err, IsNil)
	//			Check(keys, HasLen, 2)
	//		})
	//	})

})
