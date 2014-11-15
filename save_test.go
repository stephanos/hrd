package hrd

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"
)

var _ = Describe("Saver", func() {

	var (
		coll *Collection
	)

	BeforeEach(func() {
		coll = store.Coll("my-kind")
	})

	It("saves an entity", func() {
		entity := &MyModel{}

		dsPut = func(kind internal.Kind, src interface{}, completeKeys bool) ([]*internal.Key, error) {
			// TODO
			Check(completeKeys, IsFalse)
			Check(kind.Name(), Equals, "my-kind")
			return toInternalKeys(coll.NewNumKeys(42)), nil
		}

		key, err := coll.Save().Entity(entity)
		Check(err, IsNil)
		Check(key, Equals, coll.NewNumKey(42))
	})

	It("saves multiple entities", func() {
		entities := []*MyModel{&MyModel{}, &MyModel{}}

		dsPut = func(kind internal.Kind, src interface{}, completeKeys bool) ([]*internal.Key, error) {
			// TODO
			Check(completeKeys, IsFalse)
			Check(kind.Name(), Equals, "my-kind")
			return toInternalKeys(coll.NewNumKeys(1, 2)), nil
		}

		keys, err := coll.Save().Entities(entities)
		Check(err, IsNil)
		Check(keys, Equals, coll.NewNumKeys(1, 2))
	})

	It("requires complete keys", func() {
		dsPut = func(_ internal.Kind, _ interface{}, completeKeys bool) ([]*internal.Key, error) {
			Check(completeKeys, IsTrue)
			return nil, nil
		}

		coll.Save(CompleteKeys).Entity(&MyModel{})
	})
})
