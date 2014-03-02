package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("Action Options", func() {

	It("configures local cache", func() {
		opts := defaultOperationOpts()
		Check(opts.readLocalCache, IsTrue)
		Check(opts.writeLocalCache, IsTrue)

		opts1 := opts.NoLocalCacheRead()
		Check(opts1.readLocalCache, IsFalse)

		opts2 := opts.NoLocalCacheWrite()
		Check(opts2.writeLocalCache, IsFalse)

		opts3 := opts.NoLocalCache()
		Check(opts3.readLocalCache, IsFalse)
		Check(opts3.writeLocalCache, IsFalse)

		opts4 := opts.Flags(NO_CACHE)
		Check(opts4.readLocalCache, IsFalse)
		Check(opts4.writeLocalCache, IsFalse)
	})

	It("configures global cache", func() {
		opts := defaultOperationOpts()
		Check(opts.readGlobalCache, IsTrue)
		Check(opts.writeGlobalCache, IsNum, 0)

		opts1 := opts.NoGlobalCacheRead()
		Check(opts1.readGlobalCache, IsFalse)

		opts2 := opts.NoGlobalCacheWrite()
		Check(opts2.writeGlobalCache, IsNum, -1)

		opts3 := opts.NoGlobalCache()
		Check(opts3.readGlobalCache, IsFalse)
		Check(opts3.writeGlobalCache, IsNum, -1)

		opts4 := opts.Flags(NO_CACHE)
		Check(opts4.readGlobalCache, IsFalse)
		Check(opts4.writeGlobalCache, IsNum, -1)
	})
})
