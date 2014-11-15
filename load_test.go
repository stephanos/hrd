package hrd

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"
)

var _ = Describe("Loader", func() {

	var (
		coll *Collection
	)

	BeforeEach(func() {
		coll = store.Coll("my-kind")
	})

	It("loads an entity", func() {
		entity := &MyModel{}

		dsGet = func(kind internal.Kind, keys []*internal.Key, dst interface{}, _ bool, multi bool) ([]*internal.Key, error) {
			Check(multi, IsFalse)
			Check(dst, Equals, entity)
			Check(kind.Name(), Equals, "my-kind")
			Check(keys, Equals, toInternalKeys(coll.NewNumKeys(42)))
			return keys, nil
		}

		key := coll.NewNumKey(42)
		ret, _ := coll.Load().Key(key).GetOne(entity)
		Check(ret, Equals, key)
	})

	It("loads multiple entities", func() {
		entities := []*MyModel{&MyModel{}, &MyModel{}}

		dsGet = func(kind internal.Kind, keys []*internal.Key, dsts interface{}, _ bool, multi bool) ([]*internal.Key, error) {
			Check(multi, IsTrue)
			Check(dsts, Equals, entities)
			Check(kind.Name(), Equals, "my-kind")
			Check(keys, Equals, toInternalKeys(coll.NewNumKeys(1, 2)))
			return keys, nil
		}

		keys := coll.NewNumKeys(1, 2)
		ret, _ := coll.Load().Keys(keys).GetAll(entities)
		Check(ret, Equals, keys)
	})

	Context("creates single-entity loader from", func() {
		It("key", func() {
			sl := coll.Load().Key(coll.NewNumKey(42))
			Check(sl.loader.keys, Equals, coll.NewNumKeys(42))
		})

		It("numeric id", func() {
			sl := coll.Load().ID(42)
			Check(sl.loader.keys, Equals, coll.NewNumKeys(42))
		})

		It("text id", func() {
			sl := coll.Load().TextID("a")
			Check(sl.loader.keys, Equals, coll.NewTextKeys("a"))
		})
	})

	Context("creates multi-entity loader from", func() {
		It("keys", func() {
			ml := coll.Load().Keys(coll.NewNumKeys(1, 2))
			Check(ml.loader.keys, Equals, coll.NewNumKeys(1, 2))
		})

		It("numeric ids", func() {
			ml := coll.Load().IDs(1, 2)
			Check(ml.loader.keys, Equals, coll.NewNumKeys(1, 2))
		})

		It("text ids", func() {
			ml := coll.Load().TextIDs("a", "z")
			Check(ml.loader.keys, Equals, coll.NewTextKeys("a", "z"))
		})
	})
})
