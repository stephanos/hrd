package hrd

// Saver can save entities to a Collection.
type Saver struct {
	coll *Collection
	opts *operationOpts
}

// newSaver creates a new Saver for the passed store.
// The collection's options are used as default options.
func newSaver(coll *Collection) *Saver {
	return &Saver{coll, coll.opts.clone()}
}

// ==== CONFIG

// Opts applies the passed sequence of Opt to the Saver's options.
func (s *Saver) Opts(opts ...Opt) *Saver {
	s.opts = s.opts.Apply(opts...)
	return s
}

// ==== EXECUTE

// Entity saves the passed entity into the datastore.
// If its key is incomplete, the returned key will
// be a unique key generated by the datastore.
func (s *Saver) Entity(src interface{}) (*Key, error) {
	keys, err := s.put(src)
	if len(keys) == 1 {
		return keys[0], err
	}
	return nil, err
}

// Entities saves the passed entities into the datastore.
// If an entity's key is incomplete, the returned keys will
// contain a unique key generated by the datastore.
func (s *Saver) Entities(srcs interface{}) ([]*Key, error) {
	return s.put(srcs)
}

func (s *Saver) put(src interface{}) ([]*Key, error) {
	keys, err := dsPut(s.coll, src, s.opts.completeKeys)
	return newKeys(keys), err
}