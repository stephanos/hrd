package hrd

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("HRD Save", func() {

	With("default settings", func() {
		saveTests()
	})

	With("w/o global cache", func() {
		saveTests(NoGlobalCache)
	})
})

func saveTests(opts ...Opt) {

	const (
		setID = 42
	)

	var (
		genID  int64
		coll   *Collection
		loader *Loader
		entity *SimpleModel
	)

	BeforeEach(func() {
		if coll == nil {
			coll = randomColl()
		}

		loader = coll.Load(opts...)
		clearCache()
	})

	It("saves an entity without id", func() {
		entity = &SimpleModel{Text: "text", Data: []byte{1, 2, 3}}
		key, err := coll.Save().Entity(entity)

		Check(err, IsNil)
		Check(key, NotNil)

		genID = key.IntID()
		Check(genID, IsGreaterThan, 0)
		Check(entity.ID(), Equals, genID)
		Check(entity.lifecycle, Equals, []string{"before-save", "after-save"})
		Check(entity.updatedAt, Not(IsZero))
		Check(entity.createdAt, Not(IsZero))
	})

	It("saves an entity with id", func() {
		entity.SetID(setID)
		key, err := coll.Save(CompleteKeys).Entity(entity)

		Check(err, IsNil)
		Check(key, NotNil)

		Check(key.IntID(), EqualsNum, setID)
		Check(entity.ID(), EqualsNum, setID)
	})

	// ==== ERRORS

	It("does not save nil entity", func() {
		_, err := coll.Save().Entity(nil)

		Check(err, NotNil).And(Contains, "must be non-nil")
	})

	It("does not save non-struct entity", func() {
		_, err := coll.Save().Entity(setID)

		Check(err, NotNil).And(Contains, "invalid value kind").And(Contains, "int")
	})

	It("does not save entity without ID() and SetID()", func() {
		invalidMdl := &InvalidModel{}
		_, err := coll.Save().Entity(invalidMdl)

		Check(err, NotNil).And(Contains, "does not provide ID")
	})

	It("does not save complete entity without Id", func() {
		entity := &SimpleModel{}
		_, err := coll.Save().ReqKey().Entity(entity)

		Check(err, NotNil).And(Contains, "incomplete key")
	})
}
