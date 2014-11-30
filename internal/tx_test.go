package internal

import (
	. "github.com/101loops/bdd"

	ae "appengine"
)

var _ = Describe("Transact", func() {

	It("should run operations inside a transaction", func() {
		kind := randomKind()

		Transact(ctx, false, func(ctx ae.Context) error {
			key, err := Put(kind, &MyModel{}, false)
			Check(err, IsNil)
			Check(key, NotNil)

			var entity *MyModel
			keys, err := Get(kind, key, &entity, false, false)
			Check(err, IsNil)
			Check(keys[0].Synced, NotNil)

			err = DeleteKeys(kind, keys...)
			Check(err, IsNil)

			keys, err = Get(kind, key, &entity, false, false)
			Check(err, IsNil)
			Check(keys[0].Synced, IsNil)

			return nil
		})
	})
})
