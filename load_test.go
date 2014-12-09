package hrd

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"
	"github.com/101loops/hrd/internal/types"
)

var _ = Describe("Loader", func() {

	BeforeEach(func() {
		dsGet = func(_ *types.Kind, _ []*types.Key, _ interface{}, _ bool, _ bool) ([]*types.Key, error) {
			panic("unexpected call")
		}
	})

	AfterEach(func() {
		dsGet = internal.Get
	})

	It("should load an entity", func() {
		entity := &MyModel{}

		dsGet = func(kind *types.Kind, keys []*types.Key, dst interface{}, useGlobalCache bool, multi bool) ([]*types.Key, error) {
			Check(multi, IsFalse)
			Check(dst, Equals, entity)
			Check(useGlobalCache, IsTrue)
			Check(kind.Name, Equals, "my-kind")
			Check(keys, Equals, toInternalKeys(myKind.NewNumKeys(42)))
			return keys, nil
		}

		key := myKind.NewNumKey(42)
		ret, _ := myKind.Load(ctx).Key(key).GetOne(entity)
		Check(ret, Equals, key)
	})

	It("should load multiple entities", func() {
		entities := []*MyModel{&MyModel{}, &MyModel{}}

		dsGet = func(kind *types.Kind, keys []*types.Key, dsts interface{}, _ bool, multi bool) ([]*types.Key, error) {
			Check(multi, IsTrue)
			Check(dsts, Equals, entities)
			Check(kind.Name, Equals, "my-kind")
			Check(keys, Equals, toInternalKeys(myKind.NewNumKeys(1, 2)))
			return keys, nil
		}

		keys := myKind.NewNumKeys(1, 2)
		ret, _ := myKind.Load(ctx).Keys(keys).GetAll(entities)
		Check(ret, Equals, keys)
	})

	It("should be able to skip the global cache", func() {
		dsGet = func(_ *types.Kind, _ []*types.Key, _ interface{}, useGlobalCache bool, _ bool) ([]*types.Key, error) {
			Check(useGlobalCache, IsFalse)
			return nil, nil
		}

		myKind.Load(ctx).NoGlobalCache().ID(42).GetOne(nil)
	})

	Context("should create single-entity loader from", func() {
		It("key", func() {
			sl := myKind.Load(ctx).Key(myKind.NewNumKey(42))
			Check(sl.loader.keys, Equals, myKind.NewNumKeys(42))
		})

		It("numeric id", func() {
			sl := myKind.Load(ctx).ID(42)
			Check(sl.loader.keys, Equals, myKind.NewNumKeys(42))
		})

		It("text id", func() {
			sl := myKind.Load(ctx).TextID("a")
			Check(sl.loader.keys, Equals, myKind.NewTextKeys("a"))
		})
	})

	Context("should create multi-entity loader from", func() {
		It("keys", func() {
			ml := myKind.Load(ctx).Keys(myKind.NewNumKeys(1, 2))
			Check(ml.loader.keys, Equals, myKind.NewNumKeys(1, 2))
		})

		It("numeric ids", func() {
			ml := myKind.Load(ctx).IDs(1, 2)
			Check(ml.loader.keys, Equals, myKind.NewNumKeys(1, 2))
		})

		It("text ids", func() {
			ml := myKind.Load(ctx).TextIDs("a", "z")
			Check(ml.loader.keys, Equals, myKind.NewTextKeys("a", "z"))
		})
	})
})
