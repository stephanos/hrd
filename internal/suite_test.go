package internal

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/entity"
	"github.com/101loops/hrd/internal/trafo"
	"github.com/101loops/hrd/internal/types"

	"appengine/aetest"
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

	trafo.CodecSet.AddMust(MyModel{})
	trafo.CodecSet.AddMust(InvalidModel{})

	RunSpecs(t, "HRD Internal Suite")
}

func randomKind() *types.Kind {
	var n int32
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return types.NewKind(ctx, fmt.Sprintf("coll_%v", n))
}

// ==== MODELS

type InvalidModel struct{}

type MyModel struct {
	entity.NumID
	entity.CreatedTime
	entity.UpdatedTime

	Ignore    string    `datastore:"-"`
	Num       int64     `datastore:"num"`
	Data      []byte    `datastore:",index"`
	Text      string    `datastore:"html,index"`
	Time      time.Time `datastore:"timing,index,omitempty"`
	lifecycle []string
}

func (mdl *MyModel) BeforeLoad() error {
	mdl.lifecycle = append(mdl.lifecycle, "before-load")
	return nil
}

func (mdl *MyModel) AfterLoad() error {
	mdl.lifecycle = append(mdl.lifecycle, "after-load")
	return nil
}

func (mdl *MyModel) BeforeSave() error {
	mdl.lifecycle = append(mdl.lifecycle, "before-save")
	return nil
}

func (mdl *MyModel) AfterSave() error {
	mdl.lifecycle = append(mdl.lifecycle, "after-save")
	return nil
}

// ===== UTIL

func clearCache() {
	memcache.Flush(ctx)
}

func existsInDB(keys ...*types.Key) bool {
	var entity *MyModel
	keys, err := Get(types.NewKind(ctx, keys[0].Kind), keys, &entity, false, false)
	if err != nil {
		panic(err)
	}
	return len(keys) == 1 && keys[0].Synced != nil
}
