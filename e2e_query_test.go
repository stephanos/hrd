package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("HRD Query", func() {

	With("default settings", func() {
		queryTests(true)
	})

	With("w/o hybrid", func() {
		queryTests(false)
	})

	With("w/o cache", func() {
		queryTests(true, NoCache)
	})
})

func queryTests(hybrid bool, opts ...Opt) {

	var (
		coll  *Collection
		query *Query
	)

	With("simple model", func() {

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
			query = coll.Query().Hybrid(hybrid).Opts(opts...)
		})

		It("loads all entities", func() {
			// step is required because of 'eventual consistency'
			var entities []*SimpleModel
			coll.Load().IDs(-1).GetAll(&entities)
			coll.store.ClearCache()
		})

		It("counts entities", func() {
			count, err := query.GetCount()

			Check(err, IsNil)
			Check(count, EqualsNum, 4)
		})

		It("queries entity keys", func() {
			keys, _, err := query.GetKeys()

			Check(err, IsNil)
			Check(keys, HasLen, 4)
			Check(keys[0].IntID(), EqualsNum, 1)
			Check(keys[3].IntID(), EqualsNum, 4)
		})

		It("queries no entity", func() {
			var entity *SimpleModel
			err := query.Filter("html =", "nonsense").GetFirst(&entity)

			Check(err, IsNil)
			Check(entity, IsNil)
		})

		It("queries an entity", func() {
			var entity *SimpleModel
			err := query.Filter("html =", "text1").GetFirst(&entity)

			Check(err, IsNil)
			Check(entity, NotNil)
			Check(entity.id, EqualsNum, 1)
			Check(entity.Text, Equals, "text1")
		})

		It("queries an entity projection", func() {
			var entity *SimpleModel
			err := query.Project("html").GetFirst(&entity)

			Check(err, IsNil)
			Check(entity, NotNil)
			Check(entity.id, EqualsNum, 1)
			Check(entity.Data, IsEmpty)
			Check(entity.Text, Equals, "text1")
		})

		It("queries all entities", func() {
			var entities []*SimpleModel
			keys, _, err := query.GetAll(&entities)

			Check(err, IsNil)
			Check(keys, HasLen, 4)
			Check(entities, HasLen, 4)
			Check(entities[0].id, EqualsNum, 1)
			Check(keys[0].source, Equals, sourceDatastore)
		})

		It("queries filtered entities", func() {
			var entities []*SimpleModel
			keys, _, err := query.Filter("html =", "text1").GetAll(&entities)

			Check(err, IsNil)
			Check(keys, HasLen, 1)
			Check(entities, HasLen, 1)
			Check(entities[0].Text, Equals, "text1")
		})

		It("queries by ascending order", func() {
			var entities []*SimpleModel
			// TODO: var entities map[*Key]*SimpleModel
			keys, _, err := query.OrderAsc("html").GetAll(&entities)

			Check(err, IsNil)
			Check(keys, HasLen, 4)
			Check(entities, HasLen, 4)
			Check(entities[0].id, EqualsNum, 1)
			Check(entities[0].Text, Equals, "text1")
			Check(entities[3].id, EqualsNum, 4)
			Check(entities[3].Text, Equals, "text4")
		})

		It("queries by descending order", func() {
			var entities []*SimpleModel
			// TODO: var entities map[int64]*SimpleModel
			keys, _, err := query.OrderDesc("html").GetAll(&entities)

			Check(err, IsNil)
			Check(keys, HasLen, 4)
			Check(entities, HasLen, 4)
			Check(entities[0].id, EqualsNum, 4)
			Check(entities[0].Text, Equals, "text4")
			Check(entities[3].id, EqualsNum, 1)
			Check(entities[3].Text, Equals, "text1")
		})

		It("query with cursor", func() {
			var entities []*SimpleModel
			it := query.Limit(2).Run()
			keys, err := it.GetAll(&entities)

			Check(err, IsNil)
			Check(keys, HasLen, 2)

			entities = []*SimpleModel{}
			cursor, _ := it.Cursor()
			keys, _, err = query.Start(cursor).GetAll(&entities)

			Check(err, IsNil)
			Check(keys, HasLen, 2)
		})

		//		It("deletes entity", func() {
		//			err := coll.Delete().Entity(simpleMdls[0])
		//			Check(err, IsNil)
		//
		//			var entities []*SimpleModel
		//			coll.Load().IDs(1, 2, 3, 4).GetAll(&entities)
		//
		//			keys, err := query.GetKeys()
		//			Check(err, IsNil)
		//			Check(keys, HasLen, 3)
		//
		//			//count, err := query.GetCount()
		//		})
	})

	With("complex model", func() {
		// TODO
	})
}
