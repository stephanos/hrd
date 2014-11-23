package hrd

import (
	"fmt"
	"strings"

	ae "appengine"
	ds "appengine/datastore"
)

// Query represents a datastore query.
type Query struct {
	kind   *Kind
	ctx    ae.Context
	err    *error
	typeOf qryType
	logs   []string
	limit  int
	opts   *opts
	dsQry  *ds.Query
}

type qryType int

const (
	// normal query
	fullQry qryType = 1 + iota

	// only query projected fields
	projectQry

	// fetch keys first, then use batch get to only load uncached entities
	hybridQry
)

// newQuery creates a new Query for the passed kind.
// The kind's options are used as default options.
func newQuery(ctx ae.Context, kind *Kind) (ret *Query) {
	return &Query{
		ctx:    ctx,
		kind:   kind,
		limit:  -1,
		typeOf: hybridQry,
		opts:   defaultOpts(),
		dsQry:  ds.NewQuery(kind.name),
		logs:   []string{"KIND " + kind.name},
	}
}

func (qry *Query) clone() *Query {
	ret := *qry
	ret.opts = qry.opts.clone()
	if len(qry.logs) > 0 {
		ret.logs = make([]string, len(qry.logs))
		copy(ret.logs, qry.logs)
	}
	return &ret
}

// Opts returns a derivative Query with all passed-in options applied.
func (qry *Query) Opts(opts ...Opt) (ret *Query) {
	ret = qry.clone()
	ret.opts = ret.opts.Apply(opts...)
	return
}

// Hybrid returns a derivative Query which will run as a hybrid or non-hybrid
// query depending on the passed-in argument.
func (qry *Query) Hybrid(enabled bool) (ret *Query) {
	ret = qry.clone()
	if enabled {
		if ret.typeOf == fullQry {
			ret.typeOf = hybridQry
		}
	} else {
		if ret.typeOf == hybridQry {
			ret.typeOf = fullQry
		}
	}
	return ret
}

// Limit returns a derivative Query that has a limit on the number
// of results returned. A negative value means unlimited.
func (qry *Query) Limit(limit int) (ret *Query) {
	ret = qry.clone()
	if limit > 0 {
		ret.log("LIMIT %v", limit)
	} else {
		limit = -1
		ret.log("NO LIMIT")
	}
	ret.dsQry = ret.dsQry.Limit(limit)
	ret.limit = limit
	return ret
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
	ret.dsQry = ret.dsQry.Ancestor(k.ToDSKey(qry.ctx))
	return ret
}

// Project returns a derivative Query that yields only the passed fields.
// It cannot be used in a keys-only query.
func (qry *Query) Project(s ...string) (ret *Query) {
	ret = qry.clone()
	ret.log("PROJECT '%v'", strings.Join(s, "', '"))
	ret.dsQry = ret.dsQry.Project(s...)
	ret.typeOf = projectQry
	return ret
}

// EventualConsistency returns a derivative query that returns eventually
// consistent results. It only has an effect on ancestor queries.
func (qry *Query) EventualConsistency() (ret *Query) {
	ret = qry.clone()
	ret.log("EVENTUAL CONSISTENCY")
	ret.dsQry = ret.dsQry.EventualConsistency()
	return ret
}

// End returns a derivative Query with the passed end point.
func (qry *Query) End(c string) (ret *Query) {
	ret = qry.clone()
	if c != "" {
		if cursor, err := ds.DecodeCursor(c); err == nil {
			ret.log("END CURSOR")
			ret.dsQry = ret.dsQry.End(cursor)
		} else {
			err = fmt.Errorf("invalid end cursor %q: %v", c, err)
			ret.err = &err
		}
	}
	return ret
}

// Start returns a derivative Query with the passed start point.
func (qry *Query) Start(c string) (ret *Query) {
	ret = qry.clone()
	if c != "" {
		if cursor, err := ds.DecodeCursor(c); err == nil {
			ret.log("START CURSOR")
			ret.dsQry = ret.dsQry.Start(cursor)
		} else {
			err = fmt.Errorf("invalid start cursor %q: %v", c, err)
			ret.err = &err
		}
	}
	return ret
}

// Offset returns a derivative Query that has an offset of how many keys
// to skip over before returning results. A negative value is invalid.
func (qry *Query) Offset(off int) (ret *Query) {
	ret = qry.clone()
	ret.log("OFFSET %v", off)
	ret.dsQry = ret.dsQry.Offset(off)
	return
}

// OrderAsc returns a derivative Query with a field-based sort order, ascending.
// Orders are applied in the order they are added.
func (qry *Query) OrderAsc(s string) (ret *Query) {
	ret = qry.clone()
	ret.log("ORDER ASC %v", s)
	ret.dsQry = ret.dsQry.Order(s)
	return ret
}

// OrderDesc returns a derivative Query with a field-based sort order, descending.
// Orders are applied in the order they are added.
func (qry *Query) OrderDesc(s string) (ret *Query) {
	ret = qry.clone()
	ret.log("ORDER DESC %v", s)
	ret.dsQry = ret.dsQry.Order("-" + s)
	return
}

// Filter returns a derivative Query with a field-based filter.
// The filterStr argument must be a field name followed by optional space,
// followed by an operator, one of ">", "<", ">=", "<=", or "=".
// Fields are compared against the provided value using the operator.
// Multiple filters are AND'ed together.
func (qry *Query) Filter(q string, val interface{}) (ret *Query) {
	ret = qry.clone()
	ret.log("FILTER '%v %v'", q, val)
	ret.dsQry = ret.dsQry.Filter(q, val)
	return
}

// ==== EXECUTE

// GetCount returns the number of results for the query.
func (qry *Query) GetCount() (int, error) {
	qry.log("COUNT")
	qry.ctx.Infof(qry.getLog())

	if qry.err != nil {
		return 0, *qry.err
	}
	return qry.dsQry.Count(qry.ctx)
}

// GetKeys executes the query as keys-only: No entities are retrieved, just their keys.
func (qry *Query) GetKeys() ([]*Key, string, error) {
	q := qry.clone()
	q.dsQry = q.dsQry.KeysOnly()
	q.log("KEYS-ONLY")

	it := q.Run()
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
	if qry.err != nil {
		return nil, "", *qry.err
	}

	if qry.limit != 1 && qry.typeOf == hybridQry && qry.opts.useGlobalCache {
		keys, cursor, err := qry.GetKeys()
		if err == nil && len(keys) > 0 {
			keys, err = newLoader(qry.ctx, qry.kind).Keys(keys).GetAll(dsts)
		}
		return keys, cursor, err
	}

	it := qry.Run()
	keys, err := it.GetAll(dsts)
	if err != nil {
		return nil, "", err
	}

	cursor, err := it.Cursor()
	return keys, cursor, err
}

// GetFirst executes the query and writes the result's first entity
// to the passed destination.
func (qry *Query) GetFirst(dst interface{}) (err error) {
	return qry.Run().GetOne(dst)
}

// Run executes the query and returns an Iterator.
func (qry *Query) Run() *Iterator {
	qry.ctx.Infof(qry.getLog())
	return &Iterator{qry, qry.dsQry.Run(qry.ctx)}
}

func (qry *Query) log(s string, values ...interface{}) {
	qry.logs = append(qry.logs, fmt.Sprintf(s, values...))
}

func (qry *Query) getLog() string {
	return fmt.Sprintf("running query \"%v\"", strings.Join(qry.logs, " | "))
}
