package hrd

import ae "appengine"

// Loader can load entities from a kind.
type Loader struct {
	*actionContext
	keys []*Key
}

// newLoader creates a new Loader for the passed kind.
// The kind's options are used as default options.
func newLoader(ctx ae.Context, kind *Kind) *Loader {
	return &Loader{actionContext: newActionContext(ctx, kind)}
}

// NoGlobalCache prevents reading/writing entities from/to memcache.
func (l *Loader) NoGlobalCache() *Loader {
	l.opts = l.opts.Clone()
	l.opts.NoGlobalCache = true
	return l
}

// Key loads a single entity by key from the datastore.
func (l *Loader) Key(key *Key) *SingleLoader {
	l.keys = []*Key{key}
	return &SingleLoader{l}
}

// Keys load multiple entities by key from the datastore.
func (l *Loader) Keys(keys []*Key) *MultiLoader {
	l.keys = keys
	return &MultiLoader{l}
}

// ID loads a single entity by id from the datastore.
func (l *Loader) ID(id int64, parent ...*Key) *SingleLoader {
	l.keys = []*Key{l.kind.NewNumKey(id, parent...)}
	return &SingleLoader{l}
}

// TextID loads a single key by text id from the datastore.
func (l *Loader) TextID(id string, parent ...*Key) *SingleLoader {
	l.keys = []*Key{l.kind.NewTextKey(id, parent...)}
	return &SingleLoader{l}
}

// IDs load multiple keys by id from the datastore.
func (l *Loader) IDs(ids ...int64) *MultiLoader {
	l.keys = l.kind.NewNumKeys(ids...)
	return &MultiLoader{l}
}

// TextIDs load multiple keys by text id from the datastore.
func (l *Loader) TextIDs(ids ...string) *MultiLoader {
	l.keys = l.kind.NewTextKeys(ids...)
	return &MultiLoader{l}
}

// SingleLoader is a special Loader that allows to fetch exactly one entity
// from the datastore.
type SingleLoader struct {
	loader *Loader
}

// GetOne loads an entity from the datastore into the passed destination.
func (l *SingleLoader) GetOne(dst interface{}) (*Key, error) {
	keys, err := l.loader.get(dst, false)
	if len(keys) == 1 {
		return keys[0], err
	}
	return nil, err
}

func (l *Loader) get(dst interface{}, multi bool) ([]*Key, error) {
	keys, err := dsGet(l.Kind(), toInternalKeys(l.keys), dst, !l.opts.NoGlobalCache, multi)
	return importKeys(keys), err
}

// MultiLoader is a special Loader that allows to fetch multiple entities
// from the datastore.
type MultiLoader struct {
	loader *Loader
}

// TODO: func (l *Loader) GetEntity(dst interface{}) ([]*Key, error)
// TODO: func (l *Loader) GetEntities(dsts interface{}) ([]*Key, error)

// GetAll loads entities from the datastore into the passed destination.
func (l *MultiLoader) GetAll(dsts interface{}) ([]*Key, error) {
	return l.loader.get(dsts, true)
}
