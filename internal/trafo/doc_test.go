package trafo

import (
	"reflect"
	. "github.com/101loops/bdd"
)

var _ = Describe("Doc", func() {

	type MyModel struct{}
	type UnknownModel struct{}

	BeforeSuite(func() {
		CodecSet.AddMust(MyModel{})
	})

	Context("create from instance", func() {

		It("known struct", func() {
			doc, err := newDocFromInst(MyModel{})
			Check(err, IsNil)
			Check(doc, NotNil)
		})

		It("pointer to known struct", func() {
			doc, err := newDocFromInst(&MyModel{})
			Check(err, IsNil)
			Check(doc, NotNil)
		})

		It("unkown struct", func() {
			doc, err := newDocFromInst(UnknownModel{})
			Check(doc, IsNil)
			Check(err, NotNil).And(Contains, "no registered codec found for type 'trafo.UnknownModel'")
		})

		It("pointer to unknown struct", func() {
			doc, err := newDocFromInst(&UnknownModel{})
			Check(doc, IsNil)
			Check(err, NotNil).And(Contains, "no registered codec found for type 'trafo.UnknownModel'")
		})
	})

	Context("create from type", func() {

		It("pointer to known struct", func() {
			doc, err := newDocFromType(reflect.TypeOf(&MyModel{}))
			Check(err, IsNil)
			Check(doc, NotNil)
		})

		It("pointer to unknown struct", func() {
			doc, err := newDocFromType(reflect.TypeOf(&UnknownModel{}))
			Check(doc, IsNil)
			Check(err, NotNil).And(Contains, "no registered codec found for type 'trafo.UnknownModel'")
		})
	})
})
