package hrd

import . "github.com/101loops/bdd"

var _ = Describe("Store", func() {

	It("create numeric key", func() {
		key := store.NewNumKey("my-kind", 42)

		Check(key.IntID(), EqualsNum, 42)
		Check(key.Parent(), IsNil)
	})

	It("create numeric key with parent", func() {
		key := store.NewNumKey("child-kind", 42, store.NewNumKey("parent-kind", 66))

		Check(key.IntID(), EqualsNum, 42)
		Check(key.Parent(), NotNil)
		Check(key.Parent().IntID(), EqualsNum, 66)
	})

	It("create text key", func() {
		key := store.NewTextKey("my-kind", "abc")

		Check(key.StringID(), Equals, "abc")
		Check(key.Parent(), IsNil)
	})

	It("create text key with pareht", func() {
		key := store.NewTextKey("child-kind", "abc", store.NewTextKey("parent-kind", "xyz"))

		Check(key.StringID(), Equals, "abc")
		Check(key.Parent(), NotNil)
		Check(key.Parent().StringID(), Equals, "xyz")
	})
})
