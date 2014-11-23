package trafo

import (
	"github.com/101loops/hrd/entity"

	ds "appengine/datastore"
)

// Load loads the entity from datastore properties.
func (d *Doc) Load(c <-chan ds.Property) error {
	var err error
	dst := d.get()

	// event hook: before load
	if hook, ok := dst.(entity.BeforeLoader); ok {
		err = hook.BeforeLoad()
	}
	if err != nil {
		for _ = range c {
			// channel must be drained before returning ...
		}
		return err
	}

	if err = ds.LoadStruct(dst, c); err != nil {
		return err
	}

	// event hook: after load
	if hook, ok := dst.(entity.AfterLoader); ok {
		err = hook.AfterLoad()
	}

	return err
}
