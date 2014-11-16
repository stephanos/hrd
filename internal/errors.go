package internal

import (
	"fmt"

	ae "appengine"
)

func logErr(ctx ae.Context, e interface{}) error {
	err := fmt.Errorf("%v", e)
	ctx.Errorf("%v", err)
	return err
}
