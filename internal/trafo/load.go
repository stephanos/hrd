package trafo

import (
	"github.com/101loops/hrd/entity"

	ds "appengine/datastore"
)

// Load loads the entity from datastore properties.
func (doc *Doc) Load(c <-chan ds.Property) error {
	var err error
	dst := doc.get()

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

	c, err = doc.adaptProperties(c)
	if err != nil {
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

func (doc *Doc) adaptProperties(c <-chan ds.Property) (chan ds.Property, error) {
	var props []ds.Property
	for prop := range c {
		props = append(props, prop)
	}
	c2 := make(chan ds.Property, len(props))
	for _, prop := range props {
		prop2, err := doc.adaptProperty(prop)
		if err != nil {
			return nil, err
		}
		c2 <- prop2
	}
	close(c2)
	return c2, nil
}

func (doc *Doc) adaptProperty(prop ds.Property) (ds.Property, error) {
	// TODO
	return prop, nil
}
