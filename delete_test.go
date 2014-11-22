package hrd

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"
	"github.com/101loops/hrd/internal/types"
)

var _ = Describe("Deleter", func() {

	BeforeEach(func() {
		dsDelete = func(_ *types.Kind, _ interface{}, _ bool) error {
			panic("unexpected call")
		}
		dsDeleteKeys = func(_ *types.Kind, _ []*types.Key) error {
			panic("unexpected call")
		}
	})

	AfterEach(func() {
		dsDelete = internal.Delete
		dsDeleteKeys = internal.DeleteKeys
	})

	It("deletes an entity by key", func() {
		dsDeleteKeys = func(kind *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, newNumKeys(42))
			Check(kind.Name, Equals, "my-kind")
			return nil
		}

		myKind.Delete(ctx).Key(myKind.NewNumKey(42))
	})

	It("deletes multiple entities by key", func() {
		hrdKeys := []*Key{myKind.NewNumKey(1), myKind.NewNumKey(2)}

		dsDeleteKeys = func(kind *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, toInternalKeys(ctx, myKind.name, hrdKeys))
			Check(kind.Name, Equals, "my-kind")
			return nil
		}

		myKind.Delete(ctx).Keys(hrdKeys)
	})

	It("deletes an entity by numeric id", func() {
		dsDeleteKeys = func(kind *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, newNumKeys(42))
			Check(kind.Name, Equals, "my-kind")
			return nil
		}

		myKind.Delete(ctx).ID(42)
	})

	It("deletes multiple entities by numeric id", func() {
		dsDeleteKeys = func(kind *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, newNumKeys(1, 2))
			Check(kind.Name, Equals, "my-kind")
			return nil
		}

		myKind.Delete(ctx).IDs(1, 2)
	})

	It("deletes an entity by text id", func() {
		dsDeleteKeys = func(kind *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, newTextKeys("a"))
			Check(kind.Name, Equals, "my-kind")
			return nil
		}

		myKind.Delete(ctx).TextID("a")
	})

	It("deletes multiple entities by text id", func() {
		dsDeleteKeys = func(kind *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, newTextKeys("a", "z"))
			Check(kind.Name, Equals, "my-kind")
			return nil
		}

		myKind.Delete(ctx).TextIDs("a", "z")
	})

	It("deletes an entity", func() {
		entity := &MyModel{}

		dsDelete = func(kind *types.Kind, src interface{}, multi bool) error {
			Check(multi, IsFalse)
			Check(src, Equals, entity)
			Check(kind.Name, Equals, "my-kind")
			return nil
		}

		myKind.Delete(ctx).Entity(entity)
	})

	It("deletes multiple entities", func() {
		entities := []*MyModel{&MyModel{}, &MyModel{}}

		dsDelete = func(kind *types.Kind, srcs interface{}, multi bool) error {
			Check(multi, IsTrue)
			Check(srcs, Equals, entities)
			Check(kind.Name, Equals, "my-kind")
			return nil
		}

		myKind.Delete(ctx).Entities(entities)
	})
})
