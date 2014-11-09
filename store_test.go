package hrd

import . "github.com/101loops/bdd"

var _ = Describe("Store", func() {

	It("initializes and is configurable", func() {
		Check(store.opts, NotNil)
		Check(store.tx, Equals, false)
		Check(store.CreatedAt(), Not(IsZero))
	})

	It("create a collection", func() {
		coll := store.Coll("my-kind")

		Check(coll.opts, NotNil)
		Check(coll.store, Equals, store)
		Check(coll.Name(), Equals, "my-kind")
	})

	It("create numeric key", func() {
		key := store.NewNumKey("my-kind", 42)

		Check(key.IntID(), EqualsNum, 42)
		Check(key.Parent(), IsNil)
	})

	It("create numeric keys", func() {
		keys := store.NewNumKeys("my-kind", 1, 2)

		Check(keys, HasLen, 2)
		Check(keys[0].IntID(), EqualsNum, 1)
		Check(keys[1].IntID(), EqualsNum, 2)
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

	It("create text keys", func() {
		keys := store.NewTextKeys("my-kind", "abc", "xyz")

		Check(keys, HasLen, 2)
		Check(keys[0].StringID(), Equals, "abc")
		Check(keys[1].StringID(), Equals, "xyz")
	})

	It("create text key with parent", func() {
		key := store.NewTextKey("child-kind", "abc", store.NewTextKey("parent-kind", "xyz"))

		Check(key.StringID(), Equals, "abc")
		Check(key.Parent(), NotNil)
		Check(key.Parent().StringID(), Equals, "xyz")
	})
})
