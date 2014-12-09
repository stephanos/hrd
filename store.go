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
	dsGet        = internal.Get
	dsPut        = internal.Put
	dsCount      = internal.Count
	dsDelete     = internal.Delete
	dsIterate    = internal.Iterate
	dsTransact   = internal.Transact
	dsDeleteKeys = internal.DeleteKeys
)

// Store represents the App Engine datastore.
// Usually there should only be one per application.
type Store struct {
	opts      *types.Opts
	createdAt time.Time
}

// NewStore creates a new store.
func NewStore() *Store {
	store := &Store{
		createdAt: time.Now(),
		opts:      types.DefaultOpts(),
	}
	return store
}

// NoGlobalCache prevents reading/writing entities from/to memcache.
func (s *Store) NoGlobalCache() *Store {
	s.opts.NoGlobalCache = true
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
	opts *types.Opts
}

func newActionContext(ctx ae.Context, kind *Kind) *actionContext {
	return &actionContext{ctx, kind, kind.opts.Clone()}
}

func (sa *actionContext) Kind() *types.Kind {
	return types.NewKind(sa.ctx, sa.kind.name)
}
