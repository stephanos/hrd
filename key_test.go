package hrd

import (
	. "github.com/101loops/bdd"

	"appengine/datastore"
)

var _ = Describe("Key", func() {

	var coll *Collection
	var dsNumKey *datastore.Key
	var dsTextKey *datastore.Key

	BeforeEach(func() {
		collName := "coll_key"
		coll = store.Coll(collName)
		dsNumKey = datastore.NewKey(store.ctx, collName, "", 42, nil)
		dsTextKey = datastore.NewKey(store.ctx, collName, "my-key", 0, nil)
	})

	It("create numeric one", func() {
		key := newKey(dsNumKey)

		Check(key, NotNil)
		Check(key.Exists(), IsFalse)
		Check(key.IDString(), Equals, "42")
	})

	It("create textual one", func() {
		key := newKey(dsTextKey)

		Check(key, NotNil)
		Check(key.Exists(), IsFalse)
		Check(key.IDString(), Equals, "my-key")
	})

	It("create many", func() {
		keys := newKeys([]*datastore.Key{dsNumKey, dsTextKey})

		Check(keys, NotNil)
		Check(keys, HasLen, 2)
	})

	It("return memcache key", func() {
		key := newKey(dsTextKey)
		Check(toMemKey(key), Equals, "hrd:0:agtkZXZ-dGVzdGFwcHIUCxIIY29sbF9rZXkiBm15LWtleQw")

		key.version = 42
		Check(toMemKey(key), Equals, "hrd:42:agtkZXZ-dGVzdGFwcHIUCxIIY29sbF9rZXkiBm15LWtleQw")
	})
})
