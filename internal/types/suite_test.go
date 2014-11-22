package types

import (
	"testing"

	. "github.com/101loops/bdd"

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
