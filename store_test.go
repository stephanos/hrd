package hrd

import . "github.com/101loops/bdd"

var _ = Describe("Store", func() {

	It("initializes and is configurable", func() {
		Check(myStore.opts.useGlobalCache, IsTrue)
		Check(myStore.CreatedAt(), Not(IsZero))

		myStore.Opts(NoGlobalCache)
		Check(myStore.opts.useGlobalCache, IsFalse)
	})

	It("creates a kind", func() {
		newKind := myStore.Kind("new-kind")

		Check(newKind.opts, NotNil)
		Check(newKind.store, Equals, myStore)
		Check(newKind.Name(), Equals, "new-kind")
	})

	It("registers a new entity", func() {
		type MyModel1 struct{}
		err := myStore.RegisterEntity(&MyModel1{})
		Check(err, IsNil)

		err = myStore.RegisterEntity("invalid-entity")
		Check(err, NotNil)

		type MyModel2 struct{}
		myStore.RegisterEntityMust(&MyModel2{})

		Check(func() {
			myStore.RegisterEntityMust("invalid-entity")
		}, Panics)
	})
})
