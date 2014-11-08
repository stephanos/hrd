package hrd

import . "github.com/101loops/bdd"

var _ = Describe("Transactor", func() {

	It("is created and configured", func() {
		tx := store.TX()
		Check(tx, NotNil)
		Check(tx.opts.completeKeys, IsFalse)
		Check(tx.opts.txCrossGroup, IsFalse)
		Check(tx.opts.useGlobalCache, IsTrue)

		tx.GlobalCache(false)
		Check(tx.opts.useGlobalCache, IsFalse)

		tx.XG(true)
		Check(tx.opts.txCrossGroup, IsTrue)

		tx.Opts(CompleteKeys)
		Check(tx.opts.completeKeys, IsTrue)
	})
})
