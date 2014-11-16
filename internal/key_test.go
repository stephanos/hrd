package internal

import (
	"fmt"
	"time"
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/entity"

	ds "appengine/datastore"
)

var _ = Describe("Key", func() {

	It("returns string representation of Key", func() {
		key := NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)).String()
		Check(key, Equals, "Key{'my-kind', abc}")

		key = NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil)).String()
		Check(key, Equals, "Key{'my-kind', 42}")

		key = NewKey(ds.NewKey(ctx, "my-child", "", 42,
			ds.NewKey(ctx, "my-parent", "parent", 0, nil))).String()
		Check(key, Equals, "Key{'my-child', 42}[ParentKey{'my-parent', parent}]")
	})

	It("returns if it exists in the database", func() {
		key := NewKey(ds.NewKey(ctx, "my-kind", "", 1, nil))
		Check(key.Exists(), IsFalse)

		key.synced = time.Now()
		Check(key.Exists(), IsTrue)
	})

	It("returns an inner error", func() {
		key := NewKey(ds.NewKey(ctx, "my-kind", "", 1, nil))
		Check(key.Error(), IsNil)

		key.err = fmt.Errorf("an error")
		Check(key.Error(), NotNil).And(Equals, key.err)
	})

	Context("is created from", func() {

		var (
			dsKind = &dsKind{"my-kind"}
		)

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

				key, err := getKey(dsKind, &entity)
				Check(err, IsNil)
				Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil)))
			})

			It("with text id", func() {
				entity := entityWithTextID{}
				entity.SetID("abc")

				key, err := getKey(dsKind, &entity)
				Check(err, IsNil)
				Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)))
			})

			It("with numeric parent id", func() {
				entity := entityWithParentNumID{}
				entity.SetID(42)
				entity.SetParent("my-parent", 66)

				key, err := getKey(dsKind, &entity)
				Check(err, IsNil)
				Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 42,
					ds.NewKey(ctx, "my-parent", "", 66, nil))))
			})

			It("with text parent id", func() {
				entity := entityWithParentTextID{}
				entity.SetID("abc")
				entity.SetParent("my-parent", "xyz")

				key, err := getKey(dsKind, &entity)
				Check(err, IsNil)
				Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0,
					ds.NewKey(ctx, "my-parent", "xyz", 0, nil))))
			})
		})

		Context("mutliple entities", func() {

			It("in a slice", func() {
				entities := []*entityWithNumID{
					&entityWithNumID{}, &entityWithNumID{},
				}
				entities[0].SetID(1)
				entities[1].SetID(2)

				keys, err := getKeys(dsKind, entities)

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

				keys, err := getKeys(dsKind, entities)

				Check(err, IsNil)
				Check(keys, HasLen, 2)
				Check(keys[0], Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)))
				Check(keys[1], Equals, NewKey(ds.NewKey(ctx, "my-kind", "xyz", 0, nil)))
			})
		})

		It("but not an invalid entity", func() {
			entity := "invalid"
			key, err := getKey(dsKind, &entity)

			Check(key, IsNil)
			Check(err, NotNil).And(Contains, `value type "*string" does not provide ID()`)
		})

		It("but not an invalid entity collection", func() {
			invalidEntities := "invalid"
			key, err := getKeys(dsKind, &invalidEntities)

			Check(key, IsNil)
			Check(err, NotNil).And(Contains, `value must be a slice or map, but is "string"`)
		})
	})
})
