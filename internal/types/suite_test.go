package types

import (
	"testing"

	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/entity"

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

	RunSpecs(t, "HRD Types Suite")
}

// ==== MODELS

type entityWithNumID struct {
	entity.NumID
}

type entityWithTextID struct {
	entity.TextID
}

type entityWithParentNumID struct {
	entity.NumID

	parentKind string
	parentID   int64
}

func (mdl *entityWithParentNumID) Parent() (kind string, id int64) {
	return mdl.parentKind, mdl.parentID
}

func (mdl *entityWithParentNumID) SetParent(kind string, id int64) {
	mdl.parentKind = kind
	mdl.parentID = id
}

type entityWithParentTextID struct {
	entity.TextID

	parentKind string
	parentID   string
}

func (mdl *entityWithParentTextID) Parent() (kind string, id string) {
	return mdl.parentKind, mdl.parentID
}

func (mdl *entityWithParentTextID) SetParent(kind string, id string) {
	mdl.parentKind = kind
	mdl.parentID = id
}
