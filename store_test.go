package hrd

import . "github.com/101loops/bdd"

var _ = Describe("Store", func() {

	It("create numeric key", func() {
		key := store.NewNumKey("my-kind", 42)
		Check(key.IntID(), EqualsNum, 42)
	})

	It("create text key", func() {
		key := store.NewTextKey("my-kind", "abc")
		Check(key.StringID(), Equals, "abc")
	})
})
