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

		type entityWithNumID struct {
			entity.NumID
		}

		It("an entity with numeric id", func() {
			entity := entityWithNumID{}
			entity.SetID(42)

			key, err := getKey(&dsKind{"my-kind"}, &entity)
			Check(err, IsNil)
			Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil)))
		})

		type entityWithTextID struct {
			entity.TextID
		}

		It("an entity with text id", func() {
			entity := entityWithTextID{}
			entity.SetID("abc")

			key, err := getKey(&dsKind{"my-kind"}, &entity)
			Check(err, IsNil)
			Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)))
		})

		type entityWithParentNumID struct {
			entity.NumID
			entity.ParentNumID
		}

		It("an entity with numeric id", func() {
			entity := entityWithParentNumID{}
			entity.SetID(42)
			entity.SetParent("my-parent", 66)

			key, err := getKey(&dsKind{"my-kind"}, &entity)
			Check(err, IsNil)
			Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "", 42,
				ds.NewKey(ctx, "my-parent", "", 66, nil))))
		})

		type entityWithParentTextID struct {
			entity.TextID
			entity.ParentTextID
		}

		It("an entity with text id", func() {
			entity := entityWithParentTextID{}
			entity.SetID("abc")
			entity.SetParent("my-parent", "xyz")

			key, err := getKey(&dsKind{"my-kind"}, &entity)
			Check(err, IsNil)
			Check(key, Equals, NewKey(ds.NewKey(ctx, "my-kind", "abc", 0,
				ds.NewKey(ctx, "my-parent", "xyz", 0, nil))))
		})
	})
})
