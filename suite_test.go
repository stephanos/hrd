package hrd

import (
	"testing"
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal/types"

	"appengine/aetest"
)

var (
	ctx     aetest.Context
	myKind  *Kind
	myStore *Store
)

func TestSuite(t *testing.T) {
	var err error
	ctx, err = aetest.NewContext(nil)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	myStore = NewStore()
	myKind = myStore.Kind("my-kind")

	RunSpecs(t, "HRD API Suite")
}

type MyModel struct{}

func newNumKeys(ids ...int64) []*types.Key {
	return toInternalKeys(ctx, myKind.name, myKind.NewNumKeys(ids...))
}

func newTextKeys(ids ...string) []*types.Key {
	return toInternalKeys(ctx, myKind.name, myKind.NewTextKeys(ids...))
}
