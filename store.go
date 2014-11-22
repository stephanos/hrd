package hrd

import (
	"time"

	"github.com/101loops/hrd/internal"
	"github.com/101loops/hrd/internal/trafo"
	"github.com/101loops/hrd/internal/types"

	ae "appengine"
)

// datastore operations, makes it easy to stub out during testing
var (
	dsGet        = internal.DSGet
	dsPut        = internal.DSPut
	dsDelete     = internal.DSDelete
	dsTransact   = internal.DSTransact
	dsDeleteKeys = internal.DSDeleteKeys
	dsIterate    = internal.DSIterate
)

// Store represents the ds.
// Users should only need to create one store for each request.
type Store struct {

	// opts is a collection of options.
	// It controls the store's operations.
	opts *Opts

	// createdAt is the time of the store's creation
	createdAt time.Time

	// tx is whether the store is within a transaction.
	tx bool
}

// NewStore creates a new store.
func NewStore() *Store {
	store := &Store{
		createdAt: time.Now(),
		opts:      defaultOpts(),
	}
	return store
}

// Opts applies a sequence of Opt the Store's options.
func (s *Store) Opts(opts ...Opt) *Store {
	s.opts = s.opts.Apply(opts...)
	return s
}

// RegisterEntity prepares the passed-in struct type for the datastore.
// It returns an error if the type is invalid.
func (s *Store) RegisterEntity(entity interface{}) error {
	return trafo.CodecSet.Add(entity)
}

// RegisterEntityMust prepares the passed-in struct type for the datastore.
// It panics if the type is invalid.
func (s *Store) RegisterEntityMust(entity interface{}) {
	trafo.CodecSet.AddMust(entity)
}

// Kind returns a kind for the passed name.
func (s *Store) Kind(name string) *Kind {
	return newKind(s, name)
}

// TX creates a Transactor to run a transaction on the store.
func (s *Store) TX(ctx ae.Context) *Transactor {
	return newTransactor(s, ctx)
}

// CreatedAt returns the time the store was created.
func (s *Store) CreatedAt() time.Time {
	return s.createdAt
}

type actionContext struct {
	ctx  ae.Context
	kind *Kind
	opts *Opts
}

func newActionContext(ctx ae.Context, kind *Kind) *actionContext {
	return &actionContext{ctx, kind, kind.opts.clone()}
}

func (sa *actionContext) Kind() *types.Kind {
	return types.NewKind(sa.ctx, sa.kind.name)
}
