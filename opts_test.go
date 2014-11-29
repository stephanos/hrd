package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("Operation Options", func() {

	var opts *opts

	BeforeEach(func() {
		opts = defaultOpts()
	})

	It("should have default options", func() {
		Check(opts.completeKeys, IsFalse)
		Check(opts.useGlobalCache, IsTrue)
	})

	It("should set up complete key requirements", func() {
		opts = opts.Apply(CompleteKeys)
		Check(opts.completeKeys, IsTrue)
	})

	It("should set up global cache usage", func() {
		opts = opts.Apply(NoGlobalCache)
		Check(opts.useGlobalCache, IsFalse)
	})

	It("should do nothing for no parameters", func() {
		opts = opts.Apply()
	})
})
