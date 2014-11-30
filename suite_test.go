package hrd

import (
	"testing"
	. "github.com/101loops/bdd"

	"appengine/aetest"
)

var (
	ctx     aetest.Context
	myKind  *Kind
	myStore *Store
)

type MyModel struct{}

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
