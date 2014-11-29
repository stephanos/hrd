package trafo

import (
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal/fixture"
	"github.com/101loops/hrd/internal/types"
)

var _ = Describe("DocList", func() {

	var (
		kind     *types.Kind
		entities []*fixture.EntityWithNumID
	)

	type UnknownModel struct{}
	type InvalidModel struct{}

	BeforeEach(func() {
		CodecSet.AddMust(&InvalidModel{})
		kind = types.NewKind(ctx, "my-kind")

		entities = make([]*fixture.EntityWithNumID, 4)
		for i := int64(0); i < 4; i++ {
			entity := &fixture.EntityWithNumID{}
			entity.SetID(i + 1)
			entities[i] = entity
		}
	})

	Context("create readable list", func() {

		It("from struct pointer", func() {
			list, err := NewReadableDocList(kind, entities[0])
			Check(err, IsNil)
			Check(list, NotNil)
			Check(list.list, HasLen, 1)
		})

		It("from slice of struct pointers", func() {
			list, err := NewReadableDocList(kind, entities[0:2])
			Check(err, IsNil)
			Check(list, NotNil)
			Check(list.list, HasLen, 2)
		})

		It("but not from nil value", func() {
			list, err := NewReadableDocList(kind, nil)
			Check(list, IsNil)
			Check(err, NotNil).And(Contains, "value must be non-nil")
		})

		It("but not from unknown struct pointer", func() {
			list, err := NewReadableDocList(kind, &UnknownModel{})
			Check(list, IsNil)
			Check(err, NotNil).And(Contains, "no registered codec found for type 'trafo.UnknownModel'")
		})

		It("but not from invalid entity pointer", func() {
			list, err := NewReadableDocList(kind, &InvalidModel{})
			Check(list, IsNil)
			Check(err, NotNil).And(Contains, `value type "*trafo.InvalidModel" does not provide ID()`)
		})
	})
})
