package hrd

import . "github.com/101loops/bdd"

var _ = Describe("Transactor", func() {

	var (
		coll *Collection
	)

	BeforeEach(func() {
		if coll == nil {
			coll = randomColl()
		}
	})

	It("initializes and is configurable", func() {
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

	It("runs a transaction", func() {
		store.TX().Run(func(store *Store) error {
			coll := store.Coll(coll.Name())
			entity := &SimpleModel{id: 1, Text: "text1"}

			key, err := coll.Save().Entity(entity)
			Check(err, IsNil)
			Check(key, NotNil)

			err = coll.Delete().Entity(entity)
			Check(err, IsNil)

			return nil
		})
	})
})
