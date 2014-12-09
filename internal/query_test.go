package internal

import (
	"fmt"
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal/types"
)

var _ = Describe("Query", func() {

	var (
		kind *types.Kind
		//		childColl  *Kind
		query *types.Query
		//		childQuery *Query
		keys     []*types.Key
		entities []*MyModel
		entity   *MyModel
	)

	runQuery := func(dst interface{}, multi bool) ([]*types.Key, string, error) {
		it := types.NewIterator(ctx, query)
		keys, err := Iterate(it, dst, multi)
		cursor, _ := it.Cursor()
		return keys, cursor, err
	}

	BeforeEach(func() {
		entity = nil
		kind = randomKind()

		var err error
		entities = make([]*MyModel, 4)
		for i := int64(1); i < 5; i++ {
			entity := &MyModel{Num: i, Text: fmt.Sprintf("%v", i)}
			entity.SetID(i)
			entities[i-1] = entity
		}
		keys, err = Put(kind, entities, true)
		Check(err, IsNil)
		Check(keys, HasLen, 4)
		Check(keys[0].IntID, EqualsNum, 1)
		Check(keys[1].IntID, EqualsNum, 2)
		Check(keys[2].IntID, EqualsNum, 3)
		Check(keys[3].IntID, EqualsNum, 4)

		//		childEntities := []*ChildModel{
		//			&ChildModel{id: "a", parentID: 1, parentKind: coll.name},
		//			&ChildModel{id: "b", parentID: 2, parentKind: coll.name},
		//			&ChildModel{id: "c", parentID: 1, parentKind: coll.name},
		//			&ChildModel{id: "d", parentID: 2, parentKind: coll.name},
		//		}
		//		keys, err = childColl.Save(CompleteKeys).Entities(childEntities)
		//		Check(err, IsNil)
		//		Check(keys, HasLen, 4)
		//		Check(keys[0].StringID(), Equals, "a")
		//		Check(keys[1].StringID(), Equals, "b")
		//		Check(keys[2].StringID(), Equals, "c")
		//		Check(keys[3].StringID(), Equals, "d")

		// next step is required because of eventual consistency :(
		Get(kind, keys, &entities, false, true)

		query = types.NewQuery(kind.Name)
		//		childQuery = childColl.Query().Hybrid(hybrid)
	})

	It("counts entities", func() {
		count, err := Count(ctx, query)

		Check(err, IsNil)
		Check(count, EqualsNum, 4)
	})

	It("queries entity keys", func() {
		query.TypeOf = types.KeysOnlyQuery
		keys, _, err := runQuery(nil, true)

		Check(err, IsNil)
		Check(keys, HasLen, 4)
		Check(keys[0].IntID, EqualsNum, 1)
		Check(keys[3].IntID, EqualsNum, 4)
	})

	It("queries no entity", func() {
		query.Filter = append(query.Filter, types.Filter{Filter: "html =", Value: "nonsense"})
		keys, _, err := runQuery(&entity, false)

		Check(err, IsNil)
		Check(keys, HasLen, 0)
		Check(entity, IsNil)
	})

	It("queries an entity", func() {
		query.Filter = append(query.Filter, types.Filter{Filter: "num =", Value: 1})
		keys, _, err := runQuery(&entity, false)

		Check(err, IsNil)
		Check(keys, HasLen, 1)
		Check(entity, NotNil)
		Check(entity.ID(), EqualsNum, 1)
		Check(entity.Num, EqualsNum, 1)
		Check(entity.Text, Equals, "1")
		Check(entity.lifecycle, Equals, []string{"before-load", "after-load"})
	})

	It("queries an entity projection", func() {
		query.Projection = append(query.Projection, "num")
		keys, _, err := runQuery(&entity, false)

		Check(err, IsNil)
		Check(keys, HasLen, 1)
		Check(entity, NotNil)
		Check(entity.Num, Not(IsZero))
		Check(entity.Text, IsZero)
	})

	It("queries all entities", func() {
		keys, _, err := runQuery(&entities, true)

		Check(err, IsNil)
		Check(keys, HasLen, 4)
		Check(entities, HasLen, 4)
		Check(entities[0].ID(), EqualsNum, 1)
		Check(entities[1].ID(), EqualsNum, 2)
		Check(entities[2].ID(), EqualsNum, 3)
		Check(entities[3].ID(), EqualsNum, 4)
	})

	It("queries filtered entities", func() {
		query.Filter = append(query.Filter, types.Filter{Filter: "text =", Value: "4"})
		keys, _, err := runQuery(&entities, true)

		Check(err, IsNil)
		Check(keys, HasLen, 1)
		Check(entities, HasLen, 1)
		Check(entities[0].Text, Equals, "4")
	})

	It("queries by ascending order", func() {
		query.Order = append(query.Order, types.Order{FieldName: "num", Descending: false})
		keys, _, err := runQuery(&entities, true)

		Check(err, IsNil)
		Check(keys, HasLen, 4)
		Check(entities, HasLen, 4)
		Check(entities[0].ID(), EqualsNum, 1)
		Check(entities[3].ID(), EqualsNum, 4)
	})

	It("queries by descending order", func() {
		query.Order = append(query.Order, types.Order{FieldName: "num", Descending: true})
		keys, _, err := runQuery(&entities, true)

		Check(err, IsNil)
		Check(keys, HasLen, 4)
		Check(entities, HasLen, 4)
		Check(entities[0].ID(), EqualsNum, 4)
		Check(entities[3].ID(), EqualsNum, 1)
	})

	It("query with offset", func() {
		query.Offset = 2
		keys, _, err := runQuery(&entities, true)

		Check(err, IsNil)
		Check(keys, HasLen, 2)
		Check(entities, HasLen, 2)
		Check(entities[0].ID(), EqualsNum, 3)
		Check(entities[1].ID(), EqualsNum, 4)
	})

	It("query with cursor", func() {
		query.Limit = 2
		keys, cursor, err := runQuery(&entities, true)

		Check(err, IsNil)
		Check(keys, HasLen, 2)
		Check(cursor, Not(IsEmpty))

		query.Start = cursor
		keys, _, err = runQuery(&entities, true)

		Check(err, IsNil)
		Check(keys, HasLen, 2)
		Check(keys[0].IntID, EqualsNum, 3)
		Check(keys[1].IntID, EqualsNum, 4)

		query.Start = ""
		query.End = cursor
		keys, _, err = runQuery(&entities, true)

		Check(err, IsNil)
		Check(keys, HasLen, 2)
		Check(keys[0].IntID, EqualsNum, 1)
		Check(keys[1].IntID, EqualsNum, 2)
	})

	It("query with eventual consistency", func() {
		query.Eventual = true
		keys, _, err := runQuery(&entities, true)

		Check(err, IsNil)
		Check(keys, HasLen, 4)
	})
})
