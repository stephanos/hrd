package hrd

import (
	"fmt"
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"

	ae "appengine"
)

var _ = Describe("Transactor", func() {

	With("w/ cross group", func() {
		dsTransactTests(true)
	})

	With("w/o cross group", func() {
		dsTransactTests(false)
	})
})

func dsTransactTests(crossGroup bool) {

	AfterEach(func() {
		dsTransact = internal.Transact
	})

	It("runs a transaction", func() {
		txErr := fmt.Errorf("tx error")

		tx := myStore.TX(ctx).XG(crossGroup)
		Check(tx.crossGroup, Equals, crossGroup)

		dsTransact = func(ctx ae.Context, xg bool, f func(_ ae.Context) error) error {
			Check(xg, Equals, crossGroup)
			return f(ctx)
		}

		err := tx.Run(func(tx TX) error {
			Check(tx, NotNil)
			return txErr
		})
		Check(err, Contains, "tx error")
	})
}
