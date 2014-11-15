package internal

import (
	"fmt"

	"appengine"
)

func logErr(ctx appengine.Context, e interface{}) error {
	err := fmt.Errorf("%v", e)
	ctx.Errorf("%v", err)
	return err
}
