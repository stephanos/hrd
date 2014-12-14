package hrd

import (
	"github.com/101loops/hrd/internal/types"

	ae "appengine"
)

// Query represents a datastore query.
type Query struct {
	inner *types.Query
	ctx   ae.Context
	kind  *Kind
	opts  *types.Opts
}

// newQuery creates a new Query for the passed kind.
// The kind's options are used as default options.
func newQuery(ctx ae.Context, kind *Kind) (ret *Query) {
	return &Query{
		inner: types.NewQuery(kind.name),
		ctx:   ctx,
		kind:  kind,
		opts:  types.DefaultOpts(),
	}
}

func (qry *Query) clone() *Query {
	ret := *qry
	ret.opts = qry.opts.Clone()
	ret.inner = qry.inner.Clone()
	//	if len(qry.log) > 0 {
	//		ret.log = make([]string, len(qry.log))
	//		copy(ret.log, qry.log)
	//	}
	return &ret
}

// NoGlobalCache prevents reading/writing entities from/to memcache.
func (qry *Query) NoGlobalCache() (ret *Query) {
	ret = qry.clone()
	ret.opts.NoGlobalCache = true
	return
}

// Limit returns a derivative Query that has a limit on the number
// of results returned. A negative value means unlimited.
func (qry *Query) Limit(limit int) (ret *Query) {
	ret = qry.clone()
	if limit > 0 {
		//ret.log("LIMIT %v", limit)
	} else {
		limit = -1
		//ret.log("NO LIMIT")
	}
	ret.inner.Limit = limit
	return
}

// NoLimit returns a derivative Query that has no limit on the number
// of results returned.
func (qry *Query) NoLimit() (ret *Query) {
	return qry.Limit(-1)
}

// Ancestor returns a derivative Query with an ancestor filter.
// The ancestor should not be nil.
func (qry *Query) Ancestor(k *Key) (ret *Query) {
	ret = qry.clone()
	//ret.log("ANCESTOR '%v'", k.String())
	ret.inner.Ancestor = k.inner
	return
}

// Project returns a derivative Query that yields only the passed fields.
// It cannot be used in a keys-only query.
func (qry *Query) Project(fields ...string) (ret *Query) {
	ret = qry.clone()
	//ret.log("PROJECT '%v'", strings.Join(fields, "', '"))
	ret.inner.Projection = append([]string(nil), fields...)
	ret.inner.TypeOf = types.ProjectQuery
	return
}

// EventualConsistency returns a derivative query that returns eventually
// consistent results. It only has an effect on ancestor queries.
func (qry *Query) EventualConsistency() (ret *Query) {
	ret = qry.clone()
	//ret.log("EVENTUAL CONSISTENCY")
	ret.inner.Eventual = true
	return
}

// Start returns a derivative Query with the passed start point.
func (qry *Query) Start(c string) (ret *Query) {
	ret = qry.clone()
	//ret.log("START CURSOR")
	ret.inner.Start = c
	return
}

// End returns a derivative Query with the passed end point.
func (qry *Query) End(c string) (ret *Query) {
	ret = qry.clone()
	//ret.log("END CURSOR")
	ret.inner.End = c
	return
}

// Offset returns a derivative Query that has an offset of how many keys
// to skip over before returning results. A negative value is invalid.
func (qry *Query) Offset(off int) (ret *Query) {
	ret = qry.clone()
	//ret.log("OFFSET %v", off)
	ret.inner.Offset = off
	return
}

// OrderAsc returns a derivative Query with a field-based sort order, ascending.
// Orders are applied in the order they are added.
func (qry *Query) OrderAsc(s string) (ret *Query) {
	ret = qry.clone()
	//ret.log("ORDER ASC %v", s)
	ret.inner.Order = append(ret.inner.Order, types.Order{FieldName: s, Descending: false})
	return
}

// OrderDesc returns a derivative Query with a field-based sort order, descending.
// Orders are applied in the order they are added.
func (qry *Query) OrderDesc(s string) (ret *Query) {
	ret = qry.clone()
	//ret.log("ORDER DESC %v", s)
	ret.inner.Order = append(ret.inner.Order, types.Order{FieldName: s, Descending: true})
	return
}

// Distinct returns a derivative query that yields de-duplicated entities with
// respect to the set of projected fields. It is only used for projection
// queries.
func (qry *Query) Distinct() (ret *Query) {
	ret = qry.clone()
	ret.inner.Distinct = true
	return ret
}

// Filter returns a derivative Query with a field-based filter.
// The filterStr argument must be a field name followed by optional space,
// followed by an operator, one of ">", "<", ">=", "<=", or "=".
// Fields are compared against the provided value using the operator.
// Multiple filters are AND'ed together.
func (qry *Query) Filter(q string, val interface{}) (ret *Query) {
	ret = qry.clone()
	ret.inner.Filter = append(ret.inner.Filter, types.Filter{Filter: q, Value: val})
	//ret.log("FILTER '%v %v'", q, val)
	return
}

// GetCount returns the number of results for the query.
func (qry *Query) GetCount() (int, error) {
	//qry.log("COUNT")
	//qry.ctx.Infof(qry.getLog())

	return dsCount(qry.ctx, qry.inner)
}

// GetKeys executes the query as keys-only: No entities are retrieved, just their keys.
func (qry *Query) GetKeys() ([]*Key, string, error) {
	keysQry := qry.clone()
	//keysQry.log("KEYS-ONLY")
	keysQry.inner.TypeOf = types.KeysOnlyQuery

	it := keysQry.Run()
	keys, err := it.GetAll(nil)
	if err != nil {
		return nil, "", err
	}
	cursor, err := it.Cursor()
	return keys, cursor, err
}

// GetAll runs the query and writes the entities to the passed destination.
//
// Note that, if not manually disabled, queries for more than 1 item use
// a "hybrid query". This means that first a keys-only query is executed
// and then the keys are used to lookup the local and global cache as well
// as the datastore eventually. For a warm cache this usually is
// faster and cheaper than the regular query.
func (qry *Query) GetAll(dsts interface{}) ([]*Key, string, error) {
	useHybridQry := qry.inner.Limit != 1 && qry.inner.TypeOf == types.FullQuery && !qry.opts.NoGlobalCache
	if useHybridQry {
		return qry.getAllByHybrid(dsts)
	}

	it := qry.Run()
	keys, err := it.GetAll(dsts)
	if err != nil {
		return nil, "", err
	}
	cursor, err := it.Cursor()
	return keys, cursor, err
}

func (qry *Query) getAllByHybrid(dsts interface{}) ([]*Key, string, error) {
	keys, cursor, err := qry.GetKeys()
	if err == nil && len(keys) > 0 {
		keys, err = newLoader(qry.ctx, qry.kind).Keys(keys).GetAll(dsts)
	}
	return keys, cursor, err
}

// GetFirst executes the query and writes the result's first entity
// to the passed destination.
func (qry *Query) GetFirst(dst interface{}) (*Key, error) {
	return qry.Run().GetOne(dst)
}

// Run executes the query and returns an Iterator.
func (qry *Query) Run() *Iterator {
	//qry.ctx.Infof(qry.getLog())
	return newIterator(qry)
}

//func (qry *Query) log(s string, values ...interface{}) {
//	qry.log = append(qry.log, fmt.Sprintf(s, values...))
//}
//
//func (qry *Query) getLog() string {
//	return fmt.Sprintf("running query \"%v\"", strings.Join(qry.log, " | "))
//}
