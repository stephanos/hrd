package types

import (
	"reflect"
	. "github.com/101loops/bdd"
)

var _ = Describe("Query", func() {

	It("should initialize", func() {
		qry := NewQuery("my-kind")

		Check(*qry, Equals, Query{
			Filter: make([]Filter, 0),
			Order:  make([]Order, 0),
			TypeOf: FullQuery,
			kind:   "my-kind",
			Limit:  -1,
		})
	})

	It("should be cloneable", func() {
		qry := &Query{
			kind:       "my-kind",
			Order:      []Order{Order{"age", true}, Order{"name", false}},
			Filter:     []Filter{Filter{"age >", 18}},
			Projection: []string{"name", "age"},
			Ancestor:   NewKey("my-kind", "", 42, nil),
		}
		copy := qry.Clone()

		Check(qry, Equals, copy)
		Check(addrOf(qry.Order), Not(Equals), addrOf(copy.Order))
		Check(addrOf(qry.Filter), Not(Equals), addrOf(copy.Filter))
		Check(addrOf(qry.Projection), Not(Equals), addrOf(copy.Projection))
	})

	It("should convert to dastastore Query", func() {
		qry := &Query{
			kind:       "my-kind",
			Order:      []Order{Order{"age", true}, Order{"name", false}},
			Filter:     []Filter{Filter{"age >", 18}},
			Projection: []string{"name", "age"},
			Ancestor:   NewKey("my-kind", "", 42, nil),
			Start:      "my-start-cursor",
			End:        "my-end-cursor",
			TypeOf:     KeysOnlyQuery,
			Eventual:   true,
			Distinct:   true,
			Offset:     10,
			Limit:      10,
		}

		// we are just exercising each possible combination, the functionality is tested elsewhere
		dsQry := qry.ToDSQuery(ctx)
		Check(dsQry, NotNil)
	})
})

func addrOf(val interface{}) uintptr {
	return reflect.ValueOf(val).Pointer()
}
