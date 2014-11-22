package hrd

import . "github.com/101loops/bdd"

var _ = Describe("Store", func() {

	It("initializes and is configurable", func() {
		Check(myStore.opts, NotNil)
		Check(myStore.CreatedAt(), Not(IsZero))
	})

	It("creates a kind", func() {
		newKind := myStore.Kind("new-kind")

		Check(newKind.opts, NotNil)
		Check(newKind.store, Equals, myStore)
		Check(newKind.Name(), Equals, "new-kind")
	})
})
