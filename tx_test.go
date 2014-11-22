package hrd

import (
	"fmt"
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"
)

var _ = Describe("Transactor", func() {

	AfterEach(func() {
		dsGet = internal.Get
		dsPut = internal.Put
		dsDelete = internal.Delete
		dsTransact = internal.Transact
	})

	It("initializes and is configurable", func() {
		tx := store.TX(ctx)
		Check(tx, NotNil)
		Check(tx.crossGroup, IsFalse)

		tx.XG(true)
		Check(tx.crossGroup, IsTrue)
	})

	It("runs a transaction", func() {
		err := store.TX(ctx).Run(func(tx TX) error {
			Check(tx, NotNil)
			return fmt.Errorf("tx error")
		})
		Check(err, Contains, "tx error")
	})
})
