package internal

import (
	. "github.com/101loops/bdd"

	"appengine/datastore"
)

var _ = Describe("Key", func() {
	It("return string representation of Key", func() {
		sKey := NewKey(datastore.NewKey(ctx, "my-kind", "abc", 0, nil)).String()
		Check(sKey, Equals, "Key{'my-kind', abc}")

		sKey = NewKey(datastore.NewKey(ctx, "my-kind", "", 42, nil)).String()
		Check(sKey, Equals, "Key{'my-kind', 42}")

		sKey = NewKey(datastore.NewKey(ctx, "my-child", "", 42,
			datastore.NewKey(ctx, "my-parent", "parent", 0, nil))).String()
		Check(sKey, Equals, "Key{'my-child', 42}[ParentKey{'my-parent', parent}]")
	})
})
