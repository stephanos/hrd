package internal

import (
	. "github.com/101loops/bdd"

	ds "appengine/datastore"
)

var _ = Describe("Key", func() {
	It("return string representation of Key", func() {
		sKey := NewKey(ds.NewKey(ctx, "my-kind", "abc", 0, nil)).String()
		Check(sKey, Equals, "Key{'my-kind', abc}")

		sKey = NewKey(ds.NewKey(ctx, "my-kind", "", 42, nil)).String()
		Check(sKey, Equals, "Key{'my-kind', 42}")

		sKey = NewKey(ds.NewKey(ctx, "my-child", "", 42,
			ds.NewKey(ctx, "my-parent", "parent", 0, nil))).String()
		Check(sKey, Equals, "Key{'my-child', 42}[ParentKey{'my-parent', parent}]")
	})
})
