package hrd

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"
	"github.com/101loops/hrd/internal/types"
)

var _ = Describe("Deleter", func() {

	var (
		kind *Kind
	)

	BeforeEach(func() {
		kind = newKind(store, "my-kind")

		dsDelete = func(_ *types.Kind, _ interface{}, _ bool) error {
			panic("unexpected call")
		}
		dsDeleteKeys = func(_ *types.Kind, _ []*types.Key) error {
			panic("unexpected call")
		}
	})

	AfterEach(func() {
		dsDelete = internal.DSDelete
		dsDeleteKeys = internal.DSDeleteKeys
	})

	It("deletes an entity by key", func() {
		dsDeleteKeys = func(kindt *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, newNumKeys(kind, 42))
			Check(kindt.Name, Equals, "my-kind")
			return nil
		}

		kind.Delete(ctx).Key(kind.NewNumKey(42))
	})

	It("deletes multiple entities by key", func() {
		hrdKeys := []*Key{kind.NewNumKey(1), kind.NewNumKey(2)}

		dsDeleteKeys = func(kindt *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, toInternalKeys(ctx, kind.name, hrdKeys))
			Check(kindt.Name, Equals, "my-kind")
			return nil
		}

		kind.Delete(ctx).Keys(hrdKeys)
	})

	It("deletes an entity by numeric id", func() {
		dsDeleteKeys = func(kindt *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, newNumKeys(kind, 42))
			Check(kindt.Name, Equals, "my-kind")
			return nil
		}

		kind.Delete(ctx).ID(42)
	})

	It("deletes multiple entities by numeric id", func() {
		dsDeleteKeys = func(kindt *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, newNumKeys(kind, 1, 2))
			Check(kindt.Name, Equals, "my-kind")
			return nil
		}

		kind.Delete(ctx).IDs(1, 2)
	})

	It("deletes an entity by text id", func() {
		dsDeleteKeys = func(kindt *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, newTextKeys(kind, "a"))
			Check(kindt.Name, Equals, "my-kind")
			return nil
		}

		kind.Delete(ctx).TextID("a")
	})

	It("deletes multiple entities by text id", func() {
		dsDeleteKeys = func(kindt *types.Kind, keys []*types.Key) error {
			Check(keys, Equals, newTextKeys(kind, "a", "z"))
			Check(kindt.Name, Equals, "my-kind")
			return nil
		}

		kind.Delete(ctx).TextIDs("a", "z")
	})

	It("deletes an entity", func() {
		entity := &MyModel{}

		dsDelete = func(kindt *types.Kind, src interface{}, multi bool) error {
			Check(multi, IsFalse)
			Check(src, Equals, entity)
			Check(kindt.Name, Equals, "my-kind")
			return nil
		}

		kind.Delete(ctx).Entity(entity)
	})

	It("deletes multiple entities", func() {
		entities := []*MyModel{&MyModel{}, &MyModel{}}

		dsDelete = func(kindt *types.Kind, srcs interface{}, multi bool) error {
			Check(multi, IsTrue)
			Check(srcs, Equals, entities)
			Check(kindt.Name, Equals, "my-kind")
			return nil
		}

		kind.Delete(ctx).Entities(entities)
	})
})
