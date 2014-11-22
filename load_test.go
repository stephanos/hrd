package hrd

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"
	"github.com/101loops/hrd/internal/types"
)

var _ = Describe("Loader", func() {

	var (
		kind *Kind
	)

	BeforeEach(func() {
		kind = store.Kind("my-kind")
	})

	AfterEach(func() {
		dsGet = internal.DSGet
	})

	It("loads an entity", func() {
		entity := &MyModel{}

		dsGet = func(kindt *types.Kind, keys []*types.Key, dst interface{}, _ bool, multi bool) ([]*types.Key, error) {
			Check(multi, IsFalse)
			Check(dst, Equals, entity)
			Check(kindt.Name, Equals, "my-kind")
			Check(keys, Equals, newNumKeys(kind, 42))
			return keys, nil
		}

		key := kind.NewNumKey(42)
		ret, _ := kind.Load(ctx).Key(key).GetOne(entity)
		Check(ret, Equals, key)
	})

	It("loads multiple entities", func() {
		entities := []*MyModel{&MyModel{}, &MyModel{}}

		dsGet = func(kindt *types.Kind, keys []*types.Key, dsts interface{}, _ bool, multi bool) ([]*types.Key, error) {
			Check(multi, IsTrue)
			Check(dsts, Equals, entities)
			Check(kindt.Name, Equals, "my-kind")
			Check(keys, Equals, newNumKeys(kind, 1, 2))
			return keys, nil
		}

		keys := kind.NewNumKeys(1, 2)
		ret, _ := kind.Load(ctx).Keys(keys).GetAll(entities)
		Check(ret, Equals, keys)
	})

	Context("creates single-entity loader from", func() {
		It("key", func() {
			sl := kind.Load(ctx).Key(kind.NewNumKey(42))
			Check(sl.loader.keys, Equals, kind.NewNumKeys(42))
		})

		It("numeric id", func() {
			sl := kind.Load(ctx).ID(42)
			Check(sl.loader.keys, Equals, kind.NewNumKeys(42))
		})

		It("text id", func() {
			sl := kind.Load(ctx).TextID("a")
			Check(sl.loader.keys, Equals, kind.NewTextKeys("a"))
		})
	})

	Context("creates multi-entity loader from", func() {
		It("keys", func() {
			ml := kind.Load(ctx).Keys(kind.NewNumKeys(1, 2))
			Check(ml.loader.keys, Equals, kind.NewNumKeys(1, 2))
		})

		It("numeric ids", func() {
			ml := kind.Load(ctx).IDs(1, 2)
			Check(ml.loader.keys, Equals, kind.NewNumKeys(1, 2))
		})

		It("text ids", func() {
			ml := kind.Load(ctx).TextIDs("a", "z")
			Check(ml.loader.keys, Equals, kind.NewTextKeys("a", "z"))
		})
	})
})
