package hrd

import (
	"appengine/aetest"
	"appengine/memcache"
	"crypto/rand"
	"encoding/binary"
	"time"
	"fmt"
	. "github.com/101loops/bdd"
	"testing"
)

var (
	ctx   aetest.Context
	store *Store
)

func TestSuite(t *testing.T) {
	var err error
	ctx, err = aetest.NewContext(nil)
	Check(err, IsNil)
	defer ctx.Close()

	store = NewStore(ctx)

	RunSpecs(t, "HRD Suite")
}

func randomColl() *Collection {
	var n int32
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return store.Coll(fmt.Sprintf("coll_%v", n))
}

func clearCache() {
	store.ClearCache()
	memcache.Flush(ctx)
}

type InvalidModel struct{}

type SimpleModel struct {
	id        int64     `datastore:"-"`
	Num       int64     `datastore:"num"`
	Ignore    string    `datastore:"-"`
	Data      []byte    `datastore:"dat,index"`
	Text      string    `datastore:"html,index,omitempty"`
	Time      time.Time `datastore:"timing,index,omitempty"`
	lifecycle string    `datastore:"-"`
	updatedAt time.Time `datastore:"-"`
	createdAt time.Time `datastore:"-"`
}

func (self *SimpleModel) ID() int64 {
	return self.id
}

func (self *SimpleModel) SetID(id int64) {
	self.id = id
}

func (self *SimpleModel) BeforeLoad() error {
	self.lifecycle = "before-load"
	return nil
}

func (self *SimpleModel) AfterLoad() error {
	self.lifecycle = "after-load"
	return nil
}

func (self *SimpleModel) BeforeSave() error {
	self.lifecycle = "before-save"
	return nil
}

func (self *SimpleModel) AfterSave() error {
	self.lifecycle = "after-save"
	return nil
}

func (self *SimpleModel) SetCreatedAt(t time.Time) {
	self.createdAt = t
}

func (self *SimpleModel) SetUpdatedAt(t time.Time) {
	self.updatedAt = t
}

type ComplexModel struct {
	Pair Pair `datastore:"tag"`
	//PairPtr  *Pair   `datastore:"pair"`
	Pairs []Pair `datastore:"tags"`
	//PairPtrs []*Pair `datastore:"pairs"`
	lifecycle string `datastore:"-"`
}

func (self *ComplexModel) BeforeLoad() error {
	self.lifecycle = "before-load"
	return nil
}

func (self *ComplexModel) AfterLoad() error {
	self.lifecycle = "after-load"
	return nil
}

func (self *ComplexModel) BeforeSave() error {
	self.lifecycle = "before-save"
	return nil
}

func (self *ComplexModel) AfterSave() error {
	self.lifecycle = "after-save"
	return nil
}

type Pair struct {
	Key string `datastore:",index,omitempty"`
	Val string
}
