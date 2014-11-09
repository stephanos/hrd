package hrd

import . "github.com/101loops/bdd"

var _ = Describe("HRD Query", func() {

	With("default settings", func() {
		queryTests(true)
	})

	With("w/o hybrid", func() {
		queryTests(false)
	})
})

func queryTests(hybrid bool) {

	var (
		coll      *Collection
		childColl *Collection
		query     *Query
		childQuery     *Query
	)

	BeforeEach(func() {
		if coll == nil {
			coll = randomColl()
			childColl = randomColl()
		}

		entities := []*SimpleModel{
			&SimpleModel{id: 1, Text: "text1"}, &SimpleModel{id: 2, Text: "text2"},
			&SimpleModel{id: 3, Text: "text3"}, &SimpleModel{id: 4, Text: "text4"},
		}
		keys, err := coll.Save(CompleteKeys).Entities(entities)
		Check(err, IsNil)
		Check(keys, HasLen, 4)
		Check(keys[0].IntID(), EqualsNum, 1)
		Check(keys[1].IntID(), EqualsNum, 2)
		Check(keys[2].IntID(), EqualsNum, 3)
		Check(keys[3].IntID(), EqualsNum, 4)

		childEntities := []*ChildModel{
			&ChildModel{id: "a", parentID: 1, parentKind: coll.name},
			&ChildModel{id: "b", parentID: 2, parentKind: coll.name},
			&ChildModel{id: "c", parentID: 1, parentKind: coll.name},
			&ChildModel{id: "d", parentID: 2, parentKind: coll.name},
		}
		keys, err = childColl.Save(CompleteKeys).Entities(childEntities)
		Check(err, IsNil)
		Check(keys, HasLen, 4)
		Check(keys[0].StringID(), Equals, "a")
		Check(keys[1].StringID(), Equals, "b")
		Check(keys[2].StringID(), Equals, "c")
		Check(keys[3].StringID(), Equals, "d")

		clearCache()

		query = coll.Query().Hybrid(hybrid)
		childQuery = childColl.Query().Hybrid(hybrid)
	})

	It("loads all entities", func() {
		// step is required because of 'eventual consistency'
		var entities []*SimpleModel
		coll.Load().IDs(-1).GetAll(&entities)
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
		Check(entity.ID(), EqualsNum, 1)
		Check(entity.Text, Equals, "text1")
		Check(entity.lifecycle, Equals, []string{"before-load", "after-load"})
	})

	It("queries an entity projection", func() {
		var entity *SimpleModel
		err := query.Project("html").GetFirst(&entity)

		Check(err, IsNil)
		Check(entity, NotNil)
		Check(entity.ID(), EqualsNum, 1)
		Check(entity.Data, IsEmpty)
		Check(entity.Text, Equals, "text1")
	})

	It("queries all entities", func() {
		var entities []*SimpleModel
		validate := func(keys []*Key, cursor string, err error) {
			Check(err, IsNil)
			Check(keys, HasLen, 4)
			Check(cursor, Not(IsEmpty))

			Check(entities, HasLen, 4)
			Check(entities[0].id, EqualsNum, 1)
			Check(entities[1].id, EqualsNum, 2)
			Check(entities[2].id, EqualsNum, 3)
			Check(entities[3].id, EqualsNum, 4)
		}

		validate(query.GetAll(&entities))
		validate(query.NoLimit().GetAll(&entities))
		validate(query.Limit(-1).GetAll(&entities))
	})

	It("queries filtered entities", func() {
		var entities []*SimpleModel
		keys, cursor, err := query.Filter("html =", "text1").GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 1)
		Check(cursor, Not(IsEmpty))

		Check(entities, HasLen, 1)
		Check(entities[0].Text, Equals, "text1")
	})

	It("queries by ascending order", func() {
		var entities []*SimpleModel
		// TODO: var entities map[*Key]*SimpleModel
		keys, cursor, err := query.OrderAsc("html").GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 4)
		Check(cursor, Not(IsEmpty))

		Check(entities, HasLen, 4)
		Check(entities[0].ID(), EqualsNum, 1)
		Check(entities[0].Text, Equals, "text1")
		Check(entities[3].ID(), EqualsNum, 4)
		Check(entities[3].Text, Equals, "text4")
	})

	It("queries by descending order", func() {
		var entities []*SimpleModel
		// TODO: var entities map[int64]*SimpleModel
		keys, cursor, err := query.OrderDesc("html").GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 4)
		Check(cursor, Not(IsEmpty))

		Check(entities, HasLen, 4)
		Check(entities[0].ID(), EqualsNum, 4)
		Check(entities[0].Text, Equals, "text4")
		Check(entities[3].ID(), EqualsNum, 1)
		Check(entities[3].Text, Equals, "text1")
	})

	It("query with offset", func() {
		var entities []*SimpleModel
		keys, cursor, err := query.Offset(2).GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 2)
		Check(cursor, Not(IsEmpty))

		Check(entities, HasLen, 2)
		Check(entities[0].ID(), EqualsNum, 3)
		Check(entities[1].ID(), EqualsNum, 4)
	})

	It("query with cursor", func() {
		var entities []*SimpleModel
		it := query.Limit(2).Run()
		keys, err := it.GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 2)

		var cursor string
		cursor, err = it.Cursor()
		Check(err, IsNil)
		Check(cursor, Not(IsEmpty))

		// cursor: start
		entities = []*SimpleModel{}
		keys, _, err = query.Start(cursor).GetAll(&entities)
		Check(err, IsNil)
		Check(keys, HasLen, 2)
		Check(keys[0].IntID(), EqualsNum, 3)
		Check(keys[1].IntID(), EqualsNum, 4)

		// cursor: end
		entities = []*SimpleModel{}
		keys, _, err = query.End(cursor).GetAll(&entities)
		Check(err, IsNil)
		Check(keys, HasLen, 2)
		Check(keys[0].IntID(), EqualsNum, 1)
		Check(keys[1].IntID(), EqualsNum, 2)
	})

	It("query with eventual consistency", func() {
		var entities []*SimpleModel
		keys, _, err := query.EventualConsistency().GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 4)
	})

	It("query with ancestor", func() {
		var entities []*ChildModel
		keys, _, err := childQuery.Ancestor(coll.NewNumKey(1)).GetAll(&entities)

		Check(err, IsNil)
		Check(keys, HasLen, 2)
		Check(entities, HasLen, 2)
		Check(entities[0].ID(), Equals, "a")
		Check(entities[1].ID(), Equals, "c")
	})

	// ==== ERRORS

	expectError := func(q *Query, errContains string) {
		var entity *SimpleModel
		var entities []*SimpleModel

		_, err := q.GetCount()
		Check(err, Contains, errContains)

		_, _, err = q.GetKeys()
		Check(err, Contains, errContains)

		err = q.GetFirst(&entity)
		Check(err, Contains, errContains)

		_, _, err = q.GetAll(&entities)
		Check(err, Contains, errContains)
	}

	It("does not run query with invalid start cursor", func() {
		q := coll.Query().Start("nonsense")
		expectError(q, `invalid start cursor "nonsense"`)
	})

	It("does not run query with invalid end cursor", func() {
		q := coll.Query().End("nonsense")
		expectError(q, `invalid end cursor "nonsense"`)
	})
}
