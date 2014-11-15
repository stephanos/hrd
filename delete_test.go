package hrd

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"
)

var _ = Describe("Deleter", func() {

	var (
		coll *Collection
	)

	BeforeEach(func() {
		coll = store.Coll("my-kind")

		dsDelete = func(_ internal.Kind, _ interface{}, _ bool) error {
			panic("unexpected call")
		}
		dsDeleteKeys = func(_ internal.Kind, _ []*internal.Key) error {
			panic("unexpected call")
		}
	})

	It("deletes an entity by key", func() {
		dsDeleteKeys = func(kind internal.Kind, keys []*internal.Key) error {
			Check(keys, Equals, toInternalKeys(coll.NewNumKeys(42)))
			Check(kind.Name(), Equals, "my-kind")
			return nil
		}

		coll.Delete().Key(coll.NewNumKey(42))
	})

	It("deletes multiple entities by key", func() {
		hrdKeys := []*Key{coll.NewNumKey(1), coll.NewNumKey(2)}

		dsDeleteKeys = func(kind internal.Kind, keys []*internal.Key) error {
			Check(keys, Equals, toInternalKeys(hrdKeys))
			Check(kind.Name(), Equals, "my-kind")
			return nil
		}

		coll.Delete().Keys(hrdKeys)
	})

	It("deletes an entity by numeric id", func() {
		dsDeleteKeys = func(kind internal.Kind, keys []*internal.Key) error {
			Check(keys, Equals, toInternalKeys(coll.NewNumKeys(42)))
			Check(kind.Name(), Equals, "my-kind")
			return nil
		}

		coll.Delete().ID(42)
	})

	It("deletes multiple entities by numeric id", func() {
		dsDeleteKeys = func(kind internal.Kind, keys []*internal.Key) error {
			Check(keys, Equals, toInternalKeys(coll.NewNumKeys(1, 2)))
			Check(kind.Name(), Equals, "my-kind")
			return nil
		}

		coll.Delete().IDs(1, 2)
	})

	It("deletes an entity by text id", func() {
		dsDeleteKeys = func(kind internal.Kind, keys []*internal.Key) error {
			Check(keys, Equals, toInternalKeys(coll.NewTextKeys("a")))
			Check(kind.Name(), Equals, "my-kind")
			return nil
		}

		coll.Delete().TextID("a")
	})

	It("deletes multiple entities by text id", func() {
		dsDeleteKeys = func(kind internal.Kind, keys []*internal.Key) error {
			Check(keys, Equals, toInternalKeys(coll.NewTextKeys("a", "z")))
			Check(kind.Name(), Equals, "my-kind")
			return nil
		}

		coll.Delete().TextIDs("a", "z")
	})

	It("deletes an entity", func() {
		entity := &MyModel{}

		dsDelete = func(kind internal.Kind, src interface{}, multi bool) error {
			Check(multi, IsFalse)
			Check(src, Equals, entity)
			Check(kind.Name(), Equals, "my-kind")
			return nil
		}

		coll.Delete().Entity(entity)
	})

	It("deletes multiple entities", func() {
		entities := []*MyModel{&MyModel{}, &MyModel{}}

		dsDelete = func(kind internal.Kind, srcs interface{}, multi bool) error {
			Check(multi, IsTrue)
			Check(srcs, Equals, entities)
			Check(kind.Name(), Equals, "my-kind")
			return nil
		}

		coll.Delete().Entities(entities)
	})
})
