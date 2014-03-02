package hrd

type Deleter struct {
	coll *Collection
}

// Delete a single entity by key from the datastore.
func (self *Deleter) Key(key *Key) error {
	return self.delete([]*Key{key})
}

// Delete multiple entities by key from the datastore.
func (self *Deleter) Keys(keys []*Key) error {
	return self.delete(keys)
}

// Delete the provided entity.
func (self *Deleter) Entity(src interface{}) error {
	k, err := self.coll.store.getKey(self.coll.name, src)
	if err != nil {
		return err
	}
	return self.delete([]*Key{k})
}

// Delete the provided entities.
func (self *Deleter) Entities(srcs interface{}) error {
	dsKeys, err := self.coll.store.getKeys(self.coll.name, srcs)
	if err != nil {
		return err
	}
	return self.delete(dsKeys)
}

// Delete a single entity by id from the datastore.
func (self *Deleter) ID(id int64, parent ...*Key) error {
	return self.delete([]*Key{self.coll.NewNumKey(id, parent...)})
}

// Delete a single key by text id from the datastore.
func (self *Deleter) TextID(id string, parent ...*Key) error {
	return self.delete([]*Key{self.coll.NewTextKey(id, parent...)})
}

// Delete multiple keys by id from the datastore.
func (self *Deleter) IDs(ids ...int64) error {
	return self.delete(self.coll.NewNumKeys(ids...))
}

// Delete multiple keys by text id from the datastore.
func (self *Deleter) TextIDs(ids ...string) error {
	return self.delete(self.coll.NewTextKeys(ids...))
}

func (self *Deleter) delete(keys []*Key) error {
	return self.coll.store.deleteMulti(self.coll.name, keys)
}
