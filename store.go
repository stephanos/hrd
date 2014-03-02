package hrd

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"reflect"
	"time"
)

type Store struct {
	*cache
	ctx       appengine.Context
	opts      *operationOpts
	createdAt time.Time
	inTX      bool
}

func NewStore(ctx appengine.Context) *Store {
	store := &Store{
		ctx:       ctx,
		createdAt: time.Now(),
		opts:      defaultOperationOpts(),
	}
	store.cache = newStoreCache(store)
	return store
}

func (self *Store) Coll(name string) *Collection {
	return &Collection{
		store: self,
		name:  name,
		opts:  self.opts.clone(),
	}
}

func (self *Store) TX() *Transactor {
	return newTransactor(self)
}

// Clear the store's local memory cache.
func (self *Store) ClearCache() {
	self.localCache.Clear()
}

func (self *Store) NewNumKey(kind string, id int64, parent ...*Key) *Key {
	var parentKey *datastore.Key
	if len(parent) > 0 {
		parentKey = parent[0].Key
	}
	return newKey(datastore.NewKey(self.ctx, kind, "", id, parentKey))
}

func (self *Store) NewNumKeys(kind string, ids ...int64) []*Key {
	keys := make([]*Key, len(ids))
	for i, id := range ids {
		keys[i] = self.NewNumKey(kind, id)
	}
	return keys
}

func (self *Store) NewTextKey(kind string, id string, parent ...*Key) *Key {
	var parentKey *datastore.Key
	if len(parent) > 0 {
		parentKey = parent[0].Key
	}
	return newKey(datastore.NewKey(self.ctx, kind, id, 0, parentKey))
}

func (self *Store) NewTextKeys(kind string, ids ...string) []*Key {
	keys := make([]*Key, len(ids))
	for i, id := range ids {
		keys[i] = self.NewTextKey(kind, id)
	}
	return keys
}

// runTX runs f in a transaction. It calls f with a transaction context that should be
// used for all App Engine operations. Neither the local nor the global cache is touched
// during a transaction - they are updated only after a successful completion.
//
// Otherwise similar to appengine/datastore.RunInTransaction:
// https://developers.google.com/appengine/docs/go/datastore/reference#RunInTransaction
func (self *Store) runTX(f func(*Store) ([]*Key, error), opts *operationOpts) (keys []*Key, err error) {

	// execute TX
	var txStore *Store
	err = datastore.RunInTransaction(self.ctx, func(tc appengine.Context) error {
		var dsErr error
		txStore = &Store{
			ctx:  tc,
			inTX: true,
			opts: opts,
		}
		txStore.cache = newStoreCache(txStore)
		keys, dsErr = f(txStore)
		return dsErr
	}, &datastore.TransactionOptions{XG: opts.tx_cross_group})

	if err == nil {
		// update cache after successful transaction
		txStore.cache.writeTo(self.cache)
	}

	return
}

func (self *Store) getKey(kind string, src interface{}) (*Key, error) {
	var parentKey *datastore.Key
	if parented, ok := src.(numParent); ok {
		parentKey = datastore.NewKey(self.ctx, parented.ParentKind(), "", parented.Parent(), nil)
	}
	if parented, ok := src.(textParent); ok {
		parentKey = datastore.NewKey(self.ctx, parented.ParentKind(), parented.Parent(), 0, nil)
	}

	if ident, ok := src.(numIdentifier); ok {
		return newKey(datastore.NewKey(self.ctx, kind, "", ident.ID(), parentKey)), nil
	}
	if ident, ok := src.(textIdentifier); ok {
		return newKey(datastore.NewKey(self.ctx, kind, ident.ID(), 0, parentKey)), nil
	}
	return nil, fmt.Errorf("value type %q does not provide ID()", reflect.TypeOf(src))
}

func (self *Store) getKeys(kind string, src interface{}) ([]*Key, error) {
	src_val := reflect.Indirect(reflect.ValueOf(src))
	src_kind := src_val.Kind()
	if src_kind != reflect.Slice && src_kind != reflect.Map {
		return nil, fmt.Errorf("value must be a slice or map")
	}

	coll_len := src_val.Len()
	keys := make([]*Key, coll_len)

	if src_val.Kind() == reflect.Slice {
		for i := 0; i < coll_len; i++ {
			v := src_val.Index(i)
			key, err := self.getKey(kind, v.Interface())
			if err != nil {
				return nil, err
			}
			keys[i] = key
		}
		return keys, nil
	}

	for i, key := range src_val.MapKeys() {
		v := src_val.MapIndex(key)
		key, err := self.getKey(kind, v.Interface())
		if err != nil {
			return nil, err
		}
		keys[i] = key
	}
	return keys, nil
}

func (self *Store) logErr(e interface{}) error {
	err := fmt.Errorf("%v", e)
	self.ctx.Errorf("%v", err)
	return err
}

func (self *Store) logAct(verb string, prop string, keys []*Key, kind string) string {
	if len(keys) == 1 {
		id := keys[0].IdString()
		if id == "" {
			return fmt.Sprintf("%v %v %v %q", verb, "1 item", prop, kind)
		}
		id = "'" + id + "'"

		parent := ""
		if parentKey := keys[0].Parent(); parentKey != nil {
			parent = fmt.Sprintf(" (with parent '%v' from %q)", newKey(parentKey).IdString(), parentKey.Kind())
		}
		return fmt.Sprintf("%v %v %v %q%v", verb, id, prop, kind, parent)
	} else {
		return fmt.Sprintf("%v %v items %v %q", verb, len(keys), prop, kind)
	}
}
