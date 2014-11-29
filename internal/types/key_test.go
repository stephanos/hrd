package types

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal/fixture"

	ds "appengine/datastore"
)

var _ = Describe("Key", func() {

	var (
		kind *Kind
	)

	BeforeEach(func() {
		kind = NewKind(ctx, "my-kind")
	})

	It("should create a new key", func() {
		key := NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil))

		Check(key, NotNil)
		Check(key.IntID(), EqualsNum, 42)
	})

	It("should create multiple new keys", func() {
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

	It("should return string representation of Key", func() {
		str := NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)).String()
		Check(str, Equals, "Key{'my-kind', abc}")

		str = NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil)).String()
		Check(str, Equals, "Key{'my-kind', 42}")

		str = NewKey(ds.NewKey(ctx, "my-child", "", 42,
			ds.NewKey(ctx, "my-parent", "parent", 0, nil))).String()
		Check(str, Equals, "Key{'my-child', 42}[ParentKey{'my-parent', parent}]")
	})

	Context("create key from a single entity", func() {

		It("should return a new Key from numeric id", func() {
			entity := fixture.EntityWithNumID{}
			entity.SetID(42)

			key, err := GetEntityKey(kind, &entity)
			Check(err, IsNil)
			Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil)))
		})

		It("should return a new Key from text id", func() {
			entity := fixture.EntityWithTextID{}
			entity.SetID("abc")

			key, err := GetEntityKey(kind, &entity)
			Check(err, IsNil)
			Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)))
		})

		It("should return a new Key from numeric parent id", func() {
			entity := fixture.EntityWithParentNumID{}
			entity.SetID(42)
			entity.SetParent("my-parent", 66)

			key, err := GetEntityKey(kind, &entity)
			Check(err, IsNil)
			Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 42,
				ds.NewKey(ctx, "my-parent", "", 66, nil))))
		})

		It("should return a new Key from text parent id", func() {
			entity := fixture.EntityWithParentTextID{}
			entity.SetID("abc")
			entity.SetParent("my-parent", "xyz")

			key, err := GetEntityKey(kind, &entity)
			Check(err, IsNil)
			Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0,
				ds.NewKey(ctx, "my-parent", "xyz", 0, nil))))
		})

		It("should not create a Key from an invalid entity", func() {
			entity := "invalid"
			key, err := GetEntityKey(kind, &entity)

			Check(key, IsNil)
			Check(err, NotNil).And(Contains, `value type "*string" does not provide ID()`)
		})

		It("should not create a Key from an invalid entity collection", func() {
			invalidEntities := "invalid"
			key, err := GetEntitiesKeys(kind, &invalidEntities)

			Check(key, IsNil)
			Check(err, NotNil).And(Contains, `value must be a slice or map, but is "string"`)
		})
	})

	Context("mutliple entities", func() {

		It("should return a new Key from a slice", func() {
			entities := []*fixture.EntityWithNumID{
				&fixture.EntityWithNumID{}, &fixture.EntityWithNumID{},
			}
			entities[0].SetID(1)
			entities[1].SetID(2)

			keys, err := GetEntitiesKeys(kind, entities)

			Check(err, IsNil)
			Check(keys, HasLen, 2)
			Check(keys[0], Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 1, nil)))
			Check(keys[1], Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 2, nil)))
		})

		It("should return a new Key from a map", func() {
			entities := map[int]*fixture.EntityWithTextID{
				0: &fixture.EntityWithTextID{}, 1: &fixture.EntityWithTextID{},
			}
			entities[0].SetID("abc")
			entities[1].SetID("xyz")

			keys, err := GetEntitiesKeys(kind, entities)

			Check(err, IsNil)
			Check(keys, HasLen, 2)
			Check(keys[0], Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)))
			Check(keys[1], Equals, NewKey(ds.NewKey(ctx, "my-kind", "xyz", 0, nil)))
		})

		It("should not create a Key from a slice of invalid entities", func() {
			invalidEntities := []string{"invalid"}
			keys, err := GetEntitiesKeys(kind, invalidEntities)

			Check(keys, IsNil)
			Check(err, NotNil).And(Contains, `value type "string" does not provide ID()`)
		})

		It("should not create a Key from a map of invalid entities", func() {
			invalidEntities := map[int]string{0: "invalid"}
			keys, err := GetEntitiesKeys(kind, invalidEntities)

			Check(keys, IsNil)
			Check(err, NotNil).And(Contains, `value type "string" does not provide ID()`)
		})
	})
})
