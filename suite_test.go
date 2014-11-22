package hrd

import (
	"testing"
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal/types"

	"appengine/aetest"
)

var (
	ctx   aetest.Context
	store *Store
)

func TestSuite(t *testing.T) {
	var err error
	ctx, err = aetest.NewContext(nil)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	store = NewStore()

	RunSpecs(t, "HRD API Suite")
}

type MyModel struct{}

func newNumKeys(kind *Kind, ids ...int64) []*types.Key {
	return toInternalKeys(ctx, kind.name, kind.NewNumKeys(ids...))
}

func newTextKeys(kind *Kind, ids ...string) []*types.Key {
	return toInternalKeys(ctx, kind.name, kind.NewTextKeys(ids...))
}
