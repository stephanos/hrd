package trafo

import (
	"testing"
	"time"

	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/entity"
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

	CodecSet.AddMust(SimpleModel{})
	CodecSet.AddMust(ComplexModel{})

	CodecSet.AddMust(&fixture.EntityWithNumID{})
	CodecSet.AddMust(&fixture.EntityWithTextID{})
	CodecSet.AddMust(&fixture.EntityWithParentNumID{})
	CodecSet.AddMust(&fixture.EntityWithParentTextID{})

	RunSpecs(t, "HRD Trafo Suite")
}

// ==== MODELS

type SimpleModel struct {
	entity.NumID
	entity.CreatedTime
	entity.UpdatedTime

	Ignore string    `datastore:"-"`
	Num    int64     `datastore:"num"`
	Data   []byte    `datastore:",index"`
	Text   string    `datastore:"html,index"`
	Time   time.Time `datastore:"timing,index,omitempty"`

	lifecycle []string
}

func (mdl *SimpleModel) BeforeLoad() error {
	mdl.lifecycle = append(mdl.lifecycle, "before-load")
	return nil
}

func (mdl *SimpleModel) AfterLoad() error {
	mdl.lifecycle = append(mdl.lifecycle, "after-load")
	return nil
}

func (mdl *SimpleModel) BeforeSave() error {
	mdl.lifecycle = append(mdl.lifecycle, "before-save")
	return nil
}

func (mdl *SimpleModel) AfterSave() error {
	mdl.lifecycle = append(mdl.lifecycle, "after-save")
	return nil
}

type ComplexModel struct {
	Pair Pair `datastore:"tag"`
	//PairPtr  *Pair   `datastore:"pair"`
	Pairs []Pair `datastore:"tags"`
	//PairPtrs []*Pair `datastore:"pairs"`
	lifecycle string `datastore:"-"`
}

func (mdl *ComplexModel) BeforeLoad() error {
	mdl.lifecycle = "before-load"
	return nil
}

func (mdl *ComplexModel) AfterLoad() error {
	mdl.lifecycle = "after-load"
	return nil
}

func (mdl *ComplexModel) BeforeSave() error {
	mdl.lifecycle = "before-save"
	return nil
}

func (mdl *ComplexModel) AfterSave() error {
	mdl.lifecycle = "after-save"
	return nil
}

type Pair struct {
	Key string `datastore:"key,index,omitempty"`
	Val string
}
