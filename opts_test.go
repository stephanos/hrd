package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("Operation Options", func() {

	var opts *operationOpts

	BeforeEach(func() {
		opts = defaultOperationOpts()
	})

	It("has default options", func() {
		Check(opts.txCrossGroup, IsFalse)
		Check(opts.completeKeys, IsFalse)
		Check(opts.useGlobalCache, IsTrue)
		//Check(opts.useLocalCache, IsTrue)
	})

	It("configures cross-group transaction", func() {
		opts1 := opts.XG()
		Check(opts1.txCrossGroup, IsTrue)
	})

	It("configures complete key requirements", func() {
		opts1 := opts.CompleteKeys()
		Check(opts1.completeKeys, IsTrue)

		opts2 := opts1.CompleteKeys(false)
		Check(opts2.completeKeys, IsFalse)

		opts3 := opts2.CompleteKeys(true)
		Check(opts3.completeKeys, IsTrue)
	})

	It("configures global cache", func() {
		opts1 := opts.NoGlobalCache()
		Check(opts1.useGlobalCache, IsFalse)

		opts2 := opts1.GlobalCache()
		Check(opts2.useGlobalCache, IsTrue)

		opts3 := opts.Apply(NoCache)
		Check(opts3.useGlobalCache, IsFalse)
	})

})
