package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("Operation Options", func() {

	var opts *opts

	BeforeEach(func() {
		opts = defaultOpts()
	})

	It("has default options", func() {
		Check(opts.completeKeys, IsFalse)
		Check(opts.useGlobalCache, IsTrue)
	})

	It("configures complete key requirements", func() {
		opts = opts.Apply(CompleteKeys)
		Check(opts.completeKeys, IsTrue)
	})

	It("configures global cache", func() {
		opts = opts.Apply(NoGlobalCache)
		Check(opts.useGlobalCache, IsFalse)
	})
})
