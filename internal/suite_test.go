package internal

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/entity"

	ae "appengine"
	"appengine/aetest"
	ds "appengine/datastore"
	"appengine/memcache"
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
	CodecSet.AddMust(InvalidModel{})
	CodecSet.AddMust(ComplexModel{})

	RunSpecs(t, "HRD Internal Suite")
}

// ==== KIND

type dsKind struct {
	name string
}

func (k *dsKind) Name() string {
	return k.name
}

func (k *dsKind) Context() ae.Context {
	return ctx
}

func randomKind() Kind {
	var n int32
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return &dsKind{fmt.Sprintf("coll_%v", n)}
}

// ==== MODELS

type InvalidModel struct{}

type SimpleModel struct {
	entity.NumID

	Ignore    string    `datastore:"-"`
	Num       int64     `datastore:"num"`
	Data      []byte    `datastore:",index"`
	Text      string    `datastore:"html,index"`
	Time      time.Time `datastore:"timing,index,omitempty"`
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

// ===== UTIL

func clearCache() {
	memcache.Flush(ctx)
}

func existsInDB(dsKey *ds.Key) bool {
	var entity *SimpleModel
	keys, err := DSGet(&dsKind{dsKey.Kind()}, []*Key{NewKey(dsKey)}, &entity, false, false)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", keys[0].synced)
	return len(keys) == 1 && keys[0].Exists()
}
