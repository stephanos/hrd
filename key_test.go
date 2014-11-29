package hrd

import (
	"fmt"
	"time"
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal/types"
	ds "appengine/datastore"
)

var _ = Describe("Key", func() {

	It("should create numeric key", func() {
		k1 := newNumKey(myKind, 1, nil)
		Check(k1.Kind(), Equals, "my-kind")
		Check(k1.Namespace(), Equals, "")
		Check(k1.StringID(), Equals, "")
		Check(k1.IntID(), EqualsNum, 1)
		Check(k1.Incomplete(), IsFalse)
		Check(k1.Parent(), IsNil)

		k2 := newNumKey(myKind, 101, k1)
		Check(k2.Parent(), Equals, k1)
	})

	It("should create text key", func() {
		k1 := newTextKey(myKind, "abc", nil)
		Check(k1.Kind(), Equals, "my-kind")
		Check(k1.Namespace(), Equals, "")
		Check(k1.StringID(), Equals, "abc")
		Check(k1.IntID(), EqualsNum, 0)
		Check(k1.Incomplete(), IsFalse)
		Check(k1.Parent(), IsNil)

		k2 := newTextKey(myKind, "xyz", k1)
		Check(k2.Parent(), Equals, k1)
	})

	It("should create a new key", func() {
		dsKey1 := ds.NewKey(ctx, "my-kind", "abc", 0, nil)
		key1 := newKey(types.NewKey(dsKey1))
		Check(key1, Equals,
			&Key{state: &types.KeyState{}, kind: "my-kind", stringID: "abc", parent: nil})

		dsKey2 := ds.NewKey(ctx, "my-kind", "", 42, dsKey1)
		key2 := newKey(types.NewKey(dsKey2))
		Check(key2, Equals,
			&Key{state: &types.KeyState{}, kind: "my-kind", intID: 42, parent: key1})
	})

	It("should return whether it exists", func() {
		k := newTextKey(myKind, "abc", nil)
		Check(k.Exists(), IsFalse)

		k.state = nil
		Check(k.Exists(), IsFalse)

		now := time.Now()
		k.state = &types.KeyState{Synced: &now}
		Check(k.Exists(), IsTrue)
	})

	It("should return the last operation error", func() {
		k := newTextKey(myKind, "abc", nil)
		Check(k.Error(), IsNil)

		k.state = nil
		Check(k.Error(), IsNil)

		k.state = &types.KeyState{Error: fmt.Errorf("some error")}
		Check(k.Error(), NotNil)
	})

	It("should convert to datastore.Key", func() {
		k1 := newNumKey(myKind, 42, nil)
		k2 := newTextKey(myKind, "abc", k1)
		dsKey := k2.ToDSKey(ctx)

		Check(dsKey, Equals,
			ds.NewKey(ctx, "my-kind", "abc", 0, ds.NewKey(ctx, "my-kind", "", 42, nil)))
	})
})
