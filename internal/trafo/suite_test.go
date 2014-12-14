package trafo

import (
	"testing"

	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/entity/fixture"

	"appengine/aetest"
)

var (
	ctx aetest.Context
)

func TestSuite(t *testing.T) {
	var err error
	ctx, err = aetest.NewContext(nil)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	CodecSet.AddMust(&fixture.EntityWithNumID{})
	CodecSet.AddMust(&fixture.EntityWithTextID{})
	CodecSet.AddMust(&fixture.EntityWithParentNumID{})
	CodecSet.AddMust(&fixture.EntityWithParentTextID{})

	RunSpecs(t, "HRD Trafo Suite")
}

type HookEntity struct {
	A string
	B int

	beforeLoad func() error
	afterLoad  func() error
	beforeSave func() error
	afterSave  func() error
}

func (h *HookEntity) BeforeLoad() error {
	if h.beforeLoad != nil {
		return h.beforeLoad()
	}
	return nil
}

func (h *HookEntity) AfterLoad() error {
	if h.afterLoad != nil {
		return h.afterLoad()
	}
	return nil
}

func (h *HookEntity) BeforeSave() error {
	if h.beforeSave != nil {
		return h.beforeSave()
	}
	return nil
}

func (h *HookEntity) AfterSave() error {
	if h.afterSave != nil {
		return h.afterSave()
	}
	return nil
}
