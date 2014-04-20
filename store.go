package hrd

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"reflect"
	"time"
)

// Store represents the datastore.
// There should only be one store for each request.
type Store struct {
	*cache
	ctx appengine.Context

	// opts is a collection of options.
	// It controls the store's operations.
	opts *operationOpts

	// createdAt is the time of the store's creation
	createdAt time.Time

	// tx is whether the store is within a transaction.
	tx bool
}

// NewStore creates a new store for the passed App Engine context.
func NewStore(ctx appengine.Context) *Store {
	store := &Store{
		ctx:       ctx,
		createdAt: time.Now(),
		opts:      defaultOperationOpts(),
	}
	store.cache = newStoreCache(store)
	return store
}

// Coll returns a Collection for the passed name ("kind").
// The store's options are used as default options.
func (store *Store) Coll(name string) *Collection {
	return &Collection{
		store: store,
		name:  name,
		opts:  store.opts.clone(),
	}
}

// TX creates a Transactor to run a transaction on the store.
func (store *Store) TX() *Transactor {
	return newTransactor(store)
}

// ClearCache clears the store's in-memory cache.
func (store *Store) ClearCache() {
	store.localCache.Clear()
}

// CreatedAt returns the time the store was created.
func (store *Store) CreatedAt() time.Time {
	return store.createdAt
}

// NewNumKey returns a key for the passed kind and numeric ID.
// It can also receive an optional parent key.
func (store *Store) NewNumKey(kind string, id int64, parent ...*Key) *Key {
	var parentKey *datastore.Key
	if len(parent) > 0 {
		parentKey = parent[0].Key
	}
	return newKey(datastore.NewKey(store.ctx, kind, "", id, parentKey))
}

// NewNumKeys returns a sequence of key for the passed kind and
// sequence of numeric ID.
func (store *Store) NewNumKeys(kind string, ids ...int64) []*Key {
	keys := make([]*Key, len(ids))
	for i, id := range ids {
		keys[i] = store.NewNumKey(kind, id)
	}
	return keys
}

// NewTextKey returns a key for the passed kind and string ID.
// It can also receive an optional parent key.
func (store *Store) NewTextKey(kind string, id string, parent ...*Key) *Key {
	var parentKey *datastore.Key
	if len(parent) > 0 {
		parentKey = parent[0].Key
	}
	return newKey(datastore.NewKey(store.ctx, kind, id, 0, parentKey))
}

// NewTextKeys returns a sequence of keys for the passed kind and
// sequence of string ID.
func (store *Store) NewTextKeys(kind string, ids ...string) []*Key {
	keys := make([]*Key, len(ids))
	for i, id := range ids {
		keys[i] = store.NewTextKey(kind, id)
	}
	return keys
}

// runTX runs f in a transaction. It calls f with a transaction context that
// should be used for all App Engine operations. Neither the local nor the
// global cache is touched during a transaction - they are updated only after
// a successful completion.
//
// Otherwise similar to appengine/datastore.RunInTransaction:
// https://developers.google.com/appengine/docs/go/datastore/reference#RunInTransaction
func (store *Store) runTX(f func(*Store) ([]*Key, error), opts *operationOpts) (keys []*Key, err error) {

	// execute TX
	var txStore *Store
	err = datastore.RunInTransaction(store.ctx, func(tc appengine.Context) error {
		var dsErr error
		txStore = &Store{
			ctx:  tc,
			tx:   true,
			opts: opts,
		}
		txStore.cache = newStoreCache(txStore)
		keys, dsErr = f(txStore)
		return dsErr
	}, &datastore.TransactionOptions{XG: opts.txCrossGroup})

	if err == nil {
		// update cache after successful transaction
		txStore.cache.writeTo(store.cache)
	}

	return
}

func (store *Store) getKey(kind string, src interface{}) (*Key, error) {
	var parentKey *datastore.Key
	if parented, ok := src.(numParent); ok {
		parentKey = datastore.NewKey(store.ctx, parented.ParentKind(), "", parented.Parent(), nil)
	}
	if parented, ok := src.(textParent); ok {
		parentKey = datastore.NewKey(store.ctx, parented.ParentKind(), parented.Parent(), 0, nil)
	}

	if ident, ok := src.(numIdentifier); ok {
		return newKey(datastore.NewKey(store.ctx, kind, "", ident.ID(), parentKey)), nil
	}
	if ident, ok := src.(textIdentifier); ok {
		return newKey(datastore.NewKey(store.ctx, kind, ident.ID(), 0, parentKey)), nil
	}
	return nil, fmt.Errorf("value type %q does not provide ID()", reflect.TypeOf(src))
}

func (store *Store) getKeys(kind string, src interface{}) ([]*Key, error) {
	srcVal := reflect.Indirect(reflect.ValueOf(src))
	srcKind := srcVal.Kind()
	if srcKind != reflect.Slice && srcKind != reflect.Map {
		return nil, fmt.Errorf("value must be a slice or map")
	}

	collLen := srcVal.Len()
	keys := make([]*Key, collLen)

	if srcVal.Kind() == reflect.Slice {
		for i := 0; i < collLen; i++ {
			v := srcVal.Index(i)
			key, err := store.getKey(kind, v.Interface())
			if err != nil {
				return nil, err
			}
			keys[i] = key
		}
		return keys, nil
	}

	for i, key := range srcVal.MapKeys() {
		v := srcVal.MapIndex(key)
		key, err := store.getKey(kind, v.Interface())
		if err != nil {
			return nil, err
		}
		keys[i] = key
	}
	return keys, nil
}

func (store *Store) logErr(e interface{}) error {
	err := fmt.Errorf("%v", e)
	store.ctx.Errorf("%v", err)
	return err
}

func (store *Store) logAct(verb string, prop string, keys []*Key, kind string) string {
	if len(keys) == 1 {
		id := keys[0].IDString()
		if id == "" {
			return fmt.Sprintf("%v %v %v %q", verb, "1 item", prop, kind)
		}
		id = "'" + id + "'"

		parent := ""
		if parentKey := keys[0].Parent(); parentKey != nil {
			parent = fmt.Sprintf(" (with parent '%v' from %q)", newKey(parentKey).IDString(), parentKey.Kind())
		}
		return fmt.Sprintf("%v %v %v %q%v", verb, id, prop, kind, parent)
	}

	return fmt.Sprintf("%v %v items %v %q", verb, len(keys), prop, kind)
}
