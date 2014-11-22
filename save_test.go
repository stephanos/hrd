package hrd

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"
	"github.com/101loops/hrd/internal/types"
)

var _ = Describe("Saver", func() {

	var (
		kind *Kind
	)

	BeforeEach(func() {
		kind = store.Kind("my-kind")
	})

	AfterEach(func() {
		dsPut = internal.Put
	})

	It("saves an entity", func() {
		entity := &MyModel{}

		dsPut = func(kindt *types.Kind, src interface{}, completeKeys bool) ([]*types.Key, error) {
			// TODO
			Check(completeKeys, IsFalse)
			Check(kindt.Name, Equals, "my-kind")
			return newNumKeys(kind, 42), nil
		}

		key, err := kind.Save(ctx).Entity(entity)
		Check(err, IsNil)
		Check(key, Equals, kind.NewNumKey(42))
	})

	It("saves multiple entities", func() {
		entities := []*MyModel{&MyModel{}, &MyModel{}}

		dsPut = func(kindt *types.Kind, src interface{}, completeKeys bool) ([]*types.Key, error) {
			// TODO
			Check(completeKeys, IsFalse)
			Check(kindt.Name, Equals, "my-kind")
			return newNumKeys(kind, 1, 2), nil
		}

		keys, err := kind.Save(ctx).Entities(entities)
		Check(err, IsNil)
		Check(keys, Equals, kind.NewNumKeys(1, 2))
	})

	It("can require complete keys", func() {
		dsPut = func(_ *types.Kind, _ interface{}, completeKeys bool) ([]*types.Key, error) {
			Check(completeKeys, IsTrue)
			return nil, nil
		}

		kind.Save(ctx).Opts(CompleteKeys).Entity(&MyModel{})
	})
})
