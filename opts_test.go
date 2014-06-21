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
		Check(opts.readLocalCache, IsTrue)
		Check(opts.writeLocalCache, IsTrue)

		Check(opts.readGlobalCache, IsTrue)
		Check(opts.writeGlobalCache, EqualsNum, 0)

		Check(opts.txCrossGroup, IsFalse)
		Check(opts.completeKeys, IsFalse)
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

	It("configures local cache", func() {
		opts1 := opts.NoLocalCacheRead()
		Check(opts1.readLocalCache, IsFalse)

		opts2 := opts.NoLocalCacheWrite()
		Check(opts2.writeLocalCache, IsFalse)

		opts3 := opts.NoLocalCache()
		Check(opts3.readLocalCache, IsFalse)
		Check(opts3.writeLocalCache, IsFalse)

		opts4 := opts.Apply(NoCache)
		Check(opts4.readLocalCache, IsFalse)
		Check(opts4.writeLocalCache, IsFalse)

		opts5 := opts.NoCacheRead()
		Check(opts5.readLocalCache, IsFalse)

		opts6 := opts.NoCacheWrite()
		Check(opts6.writeLocalCache, IsFalse)
	})

	It("configures global cache", func() {
		opts1 := opts.NoGlobalCacheRead()
		Check(opts1.readGlobalCache, IsFalse)

		opts2 := opts.NoGlobalCacheWrite()
		Check(opts2.writeGlobalCache, EqualsNum, -1)

		opts3 := opts.NoGlobalCache()
		Check(opts3.readGlobalCache, IsFalse)
		Check(opts3.writeGlobalCache, EqualsNum, -1)

		opts4 := opts.Apply(NoCache)
		Check(opts4.readGlobalCache, IsFalse)
		Check(opts4.writeGlobalCache, EqualsNum, -1)

		opts5 := opts.NoCacheRead()
		Check(opts5.readGlobalCache, IsFalse)

		opts6 := opts.NoCacheWrite()
		Check(opts6.writeGlobalCache, EqualsNum, -1)
	})

})
