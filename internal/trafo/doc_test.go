package trafo

import . "github.com/101loops/bdd"

var _ = Describe("Doc", func() {

	It("is created from instance", func() {
		doc, err := newDocFromInst(&SimpleModel{
			Num:  42,
			Text: "html",
			Data: []byte("byte"),
		})

		Check(err, IsNil)
		Check(doc, NotNil)
	})
})
