package hrd

import (
	"testing"
	. "github.com/101loops/bdd"

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

	store = NewStore(ctx)

	RunSpecs(t, "HRD API Suite")
}

type MyModel struct{}
