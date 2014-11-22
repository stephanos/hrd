package hrd

import . "github.com/101loops/bdd"

var _ = Describe("Store", func() {

	It("initializes and is configurable", func() {
		Check(store.opts, NotNil)
		Check(store.tx, Equals, false)
		Check(store.CreatedAt(), Not(IsZero))
	})

	It("creates a kind", func() {
		kind := store.Kind("my-kind")

		Check(kind.opts, NotNil)
		Check(kind.store, Equals, store)
		Check(kind.Name(), Equals, "my-kind")
	})
})
