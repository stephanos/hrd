package trafo

import (
	ae "appengine"
	ds "appengine/datastore"
)

// DocsPipe can load/save datastore entries from/to entities.
type DocsPipe struct {
	ctx  ae.Context
	Docs []*Doc
}

// Properties returns a sequence of datastore.PropertyLoadSaver.
func (p *DocsPipe) Properties() []ds.PropertyLoadSaver {
	pipes := make([]ds.PropertyLoadSaver, len(p.Docs))
	for i, doc := range p.Docs {
		pipes[i] = &docPipe{p.ctx, doc}
	}
	return pipes
}

// docPipe can load/save datastore properties from/to an entity.
type docPipe struct {
	ctx ae.Context
	doc *Doc
}

func (p *docPipe) Load(c <-chan ds.Property) error {
	return p.doc.Load(c)
}

func (p *docPipe) Save(c chan<- ds.Property) error {
	defer close(c)

	props, err := p.doc.Save(p.ctx)
	if err != nil {
		return err
	}

	for _, prop := range props {
		c <- *prop
	}

	return nil
}
