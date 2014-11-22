package types

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/entity"

	ds "appengine/datastore"
)

var _ = Describe("Key", func() {

	var (
		kind *Kind
	)

	BeforeEach(func() {
		kind = NewKind(ctx, "my-kind")
	})

	It("creates a new key", func() {
		key := NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil))

		Check(key, NotNil)
		Check(key.IntID(), EqualsNum, 42)
	})

	It("creates multiple new keys", func() {
		dsKeys := []*ds.Key{
			ds.NewKey(ctx, kind.Name, "", 1, nil), ds.NewKey(ctx, kind.Name, "", 2, nil),
		}
		keys := NewKeys(dsKeys...)

		Check(keys, HasLen, 2)
		Check(keys[0].IntID(), EqualsNum, 1)
		Check(keys[0].Kind(), Equals, kind.Name)
		Check(keys[1].IntID(), EqualsNum, 2)
		Check(keys[1].Kind(), Equals, kind.Name)
	})

	It("returns string representation of Key", func() {
		str := NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)).String()
		Check(str, Equals, "Key{'my-kind', abc}")

		str = NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil)).String()
		Check(str, Equals, "Key{'my-kind', 42}")

		str = NewKey(ds.NewKey(ctx, "my-child", "", 42,
			ds.NewKey(ctx, "my-parent", "parent", 0, nil))).String()
		Check(str, Equals, "Key{'my-child', 42}[ParentKey{'my-parent', parent}]")
	})

	Context("is created from", func() {

		type entityWithNumID struct {
			entity.NumID
		}

		type entityWithTextID struct {
			entity.TextID
		}

		type entityWithParentNumID struct {
			entity.NumID
			entity.ParentNumID
		}

		type entityWithParentTextID struct {
			entity.TextID
			entity.ParentTextID
		}

		Context("a single entity", func() {
			It("with numeric id", func() {
				entity := entityWithNumID{}
				entity.SetID(42)

				key, err := GetEntityKey(kind, &entity)
				Check(err, IsNil)
				Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil)))
			})

			It("with text id", func() {
				entity := entityWithTextID{}
				entity.SetID("abc")

				key, err := GetEntityKey(kind, &entity)
				Check(err, IsNil)
				Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)))
			})

			It("with numeric parent id", func() {
				entity := entityWithParentNumID{}
				entity.SetID(42)
				entity.SetParent("my-parent", 66)

				key, err := GetEntityKey(kind, &entity)
				Check(err, IsNil)
				Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 42,
					ds.NewKey(ctx, "my-parent", "", 66, nil))))
			})

			It("with text parent id", func() {
				entity := entityWithParentTextID{}
				entity.SetID("abc")
				entity.SetParent("my-parent", "xyz")

				key, err := GetEntityKey(kind, &entity)
				Check(err, IsNil)
				Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0,
					ds.NewKey(ctx, "my-parent", "xyz", 0, nil))))
			})

			It("but not an invalid entity", func() {
				entity := "invalid"
				key, err := GetEntityKey(kind, &entity)

				Check(key, IsNil)
				Check(err, NotNil).And(Contains, `value type "*string" does not provide ID()`)
			})

			It("but not an invalid entity collection", func() {
				invalidEntities := "invalid"
				key, err := GetEntitiesKeys(kind, &invalidEntities)

				Check(key, IsNil)
				Check(err, NotNil).And(Contains, `value must be a slice or map, but is "string"`)
			})

		})

		Context("mutliple entities", func() {

			It("in a slice", func() {
				entities := []*entityWithNumID{
					&entityWithNumID{}, &entityWithNumID{},
				}
				entities[0].SetID(1)
				entities[1].SetID(2)

				keys, err := GetEntitiesKeys(kind, entities)

				Check(err, IsNil)
				Check(keys, HasLen, 2)
				Check(keys[0], Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 1, nil)))
				Check(keys[1], Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 2, nil)))
			})

			It("in a map", func() {
				entities := map[int]*entityWithTextID{
					0: &entityWithTextID{}, 1: &entityWithTextID{},
				}
				entities[0].SetID("abc")
				entities[1].SetID("xyz")

				keys, err := GetEntitiesKeys(kind, entities)

				Check(err, IsNil)
				Check(keys, HasLen, 2)
				Check(keys[0], Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)))
				Check(keys[1], Equals, NewKey(ds.NewKey(ctx, "my-kind", "xyz", 0, nil)))
			})

			It("but not a slice of invalid entities", func() {
				invalidEntities := []string{"invalid"}
				keys, err := GetEntitiesKeys(kind, invalidEntities)

				Check(keys, IsNil)
				Check(err, NotNil).And(Contains, `value type "string" does not provide ID()`)
			})

			It("but not a map of invalid entities", func() {
				invalidEntities := map[int]string{0: "invalid"}
				keys, err := GetEntitiesKeys(kind, invalidEntities)

				Check(keys, IsNil)
				Check(err, NotNil).And(Contains, `value type "string" does not provide ID()`)
			})
		})
	})
})
