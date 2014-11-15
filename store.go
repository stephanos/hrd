package hrd

import (
	"time"

	"github.com/101loops/hrd/internal"

	"appengine"
	"appengine/datastore"
)

// datastore operations, makes it easy to stub out during testing
var (
	dsGet        = internal.DSGet
	dsPut        = internal.DSPut
	dsDelete     = internal.DSDelete
	dsDeleteKeys = internal.DSDeleteKeys
	dsIterate    = internal.DSIterate
)

// Store represents the datastore.
// Users should only need to create one store for each request.
type Store struct {
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
	return store
}

// RegisterEntity prepares the passed-in struct type for the datastore.
// It returns an error if the type is invalid.
func (store *Store) RegisterEntity(entity interface{}) error {
	return internal.CodecSet.Add(entity)
}

// RegisterEntityMust prepares the passed-in struct type for the datastore.
// It panics if the type is invalid.
func (store *Store) RegisterEntityMust(entity interface{}) {
	internal.CodecSet.AddMust(entity)
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

// CreatedAt returns the time the store was created.
func (store *Store) CreatedAt() time.Time {
	return store.createdAt
}

// NewNumKey returns a key for the passed kind and numeric ID.
// It can also receive an optional parent key.
func (store *Store) NewNumKey(kind string, id int64, parent ...*Key) *Key {
	var parentKey *datastore.Key
	if len(parent) > 0 {
		parentKey = parent[0].Key.Key
	}
	return newKey(internal.NewKey(datastore.NewKey(store.ctx, kind, "", id, parentKey)))
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
		parentKey = parent[0].Key.Key
	}
	return newKey(internal.NewKey(datastore.NewKey(store.ctx, kind, id, 0, parentKey)))
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
// should be used for all App Engine operations.
//
// Otherwise similar to appengine/datastore.RunInTransaction:
// https://developers.google.com/appengine/docs/go/datastore/reference#RunInTransaction
func (store *Store) runTX(f func(*Store) error, opts *operationOpts) error {
	return datastore.RunInTransaction(store.ctx, func(tc appengine.Context) error {
		var dsErr error
		txStore := &Store{
			ctx:  tc,
			tx:   true,
			opts: opts,
		}
		dsErr = f(txStore)
		return dsErr
	}, &datastore.TransactionOptions{XG: opts.txCrossGroup})
}
