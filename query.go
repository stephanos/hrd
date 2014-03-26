package hrd

import (
	"appengine/datastore"
	"fmt"
	"strings"
	"time"
)

type Query struct {
	err    *error
	coll   *Collection
	typeOf qryType
	logs   []string
	limit  int
	opts   *operationOpts
	qry    *datastore.Query
}

type qryType int

const (
	// normal query
	FullQry qryType = 1 + iota

	// only query projected fields
	ProjQry

	// fetch keys first, then use batch get to only load uncached entities
	HybridQry
)

func newQuery(coll *Collection) (ret *Query) {
	return &Query{
		coll:   coll,
		limit:  -1,
		typeOf: HybridQry,
		opts:   defaultOperationOpts(),
		qry:    datastore.NewQuery(coll.name),
		logs:   []string{"KIND " + coll.name},
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

func (qry *Query) Flags(flags ...Flag) (ret *Query) {
	ret = qry.clone()
	ret.opts = ret.opts.Flags(flags...)
	return
}

func (qry *Query) NoHybrid() *Query {
	return qry.Hybrid(false)
}

func (qry *Query) Hybrid(enabled bool) (ret *Query) {
	ret = qry.clone()
	if enabled {
		if ret.typeOf == FullQry {
			ret.typeOf = HybridQry
		}
	} else {
		if ret.typeOf == HybridQry {
			ret.typeOf = FullQry
		}
	}
	return ret
}

func (qry *Query) Limit(limit int) (ret *Query) {
	ret = qry.clone()
	if limit > 0 {
		ret.log("LIMIT %v", limit)
	} else {
		limit = -1
		ret.log("NO LIMIT")
	}
	ret.qry = ret.qry.Limit(limit)
	ret.limit = limit
	return ret
}

func (qry *Query) NoLimit() (ret *Query) {
	return qry.Limit(-1)
}

func (qry *Query) Ancestor(k *Key) (ret *Query) {
	ret = qry.clone()
	ret.log("ANCESTOR '%v'", k.IdString())
	ret.qry = ret.qry.Ancestor(k.Key)
	return ret
}

func (qry *Query) Project(s ...string) (ret *Query) {
	ret = qry.clone()
	ret.log("PROJECT '%v'", strings.Join(s, "', '"))
	ret.qry = ret.qry.Project(s...)
	ret.typeOf = ProjQry
	return ret
}

func (qry *Query) End(c string) (ret *Query) {
	ret = qry.clone()
	if c != "" {
		if cursor, err := datastore.DecodeCursor(c); err == nil {
			ret.log("END CURSOR")
			ret.qry = ret.qry.End(cursor)
		} else {
			err = fmt.Errorf("invalid end cursor (%v)", err)
			ret.err = &err
		}
	}
	return ret
}

func (qry *Query) Start(c string) (ret *Query) {
	ret = qry.clone()
	if c != "" {
		if cursor, err := datastore.DecodeCursor(c); err == nil {
			ret.log("START CURSOR")
			ret.qry = ret.qry.Start(cursor)
		} else {
			err = fmt.Errorf("invalid start cursor (%v)", err)
			ret.err = &err
		}
	}
	return ret
}

func (qry *Query) Offset(off int) (ret *Query) {
	ret = qry.clone()
	ret.log("OFFSET %v", off)
	ret.qry = ret.qry.Offset(off)
	return
}

func (qry *Query) OrderAsc(s string) (ret *Query) {
	ret = qry.clone()
	ret.log("ORDER ASC %v", s)
	ret.qry = ret.qry.Order(s)
	return ret
}

func (qry *Query) OrderDesc(s string) (ret *Query) {
	ret = qry.clone()
	ret.log("ORDER DESC %v", s)
	ret.qry = ret.qry.Order("-" + s)
	return
}

func (qry *Query) Filter(q string, val interface{}) (ret *Query) {
	ret = qry.clone()
	ret.log("FILTER '%v %v'", q, val)
	ret.qry = ret.qry.Filter(q, val)
	return
}

// ==== CACHE

func (qry *Query) NoCache() (ret *Query) {
	return qry.NoLocalCache().NoGlobalCache()
}

func (qry *Query) NoLocalCache() (ret *Query) {
	return qry.NoLocalCacheWrite().NoLocalCacheRead()
}

func (qry *Query) NoGlobalCache() (ret *Query) {
	return qry.NoGlobalCacheWrite().NoGlobalCacheRead()
}

func (qry *Query) CacheExpire(exp time.Duration) (ret *Query) {
	q := qry.clone()
	q.opts = q.opts.CacheExpire(exp)
	return q
}

func (qry *Query) NoCacheRead() (ret *Query) {
	return qry.NoGlobalCacheRead().NoLocalCacheRead()
}

func (qry *Query) NoLocalCacheRead() (ret *Query) {
	q := qry.clone()
	q.opts = q.opts.NoLocalCacheRead()
	return q
}

func (qry *Query) NoGlobalCacheRead() (ret *Query) {
	q := qry.clone()
	q.opts = q.opts.NoGlobalCacheRead()
	return q
}

func (qry *Query) NoCacheWrite() (ret *Query) {
	return qry.NoGlobalCacheWrite().NoLocalCacheWrite()
}

func (qry *Query) NoLocalCacheWrite() (ret *Query) {
	q := qry.clone()
	q.opts = q.opts.NoLocalCacheWrite()
	return q
}

func (qry *Query) NoGlobalCacheWrite() (ret *Query) {
	q := qry.clone()
	q.opts = q.opts.NoGlobalCacheWrite()
	return q
}

// ==== EXECUTE

func (qry *Query) GetCount() (int, error) {
	qry.log("COUNT")
	qry.coll.store.ctx.Infof(qry.getLog())

	if qry.err != nil {
		return 0, *qry.err
	}
	return qry.qry.Count(qry.coll.store.ctx)
}

// Runs the query as keys-only: No entities are retrieved, just their keys.
func (qry *Query) GetKeys() ([]*Key, string, error) {
	q := qry.clone()
	q.qry = q.qry.KeysOnly()
	q.log("KEYS-ONLY")

	it := q.Run()
	keys, err := it.GetAll(nil)
	if err != nil {
		return nil, "", err
	}
	cursor, err := it.Cursor()
	return keys, cursor, err
}

// Runs the query and writes the entities to the passed destination.
//
// Note that, if not manually disabled, queries for more than 1 item use a "hybrid query".
// This means that first a keys-only query is executed and then the keys are used to lookup the
// local and global cache as well as the datastore eventually. For a warm cache this usually is
// faster and cheaper than the regular query.
func (qry *Query) GetAll(dsts interface{}) ([]*Key, string, error) {
	if qry.err != nil {
		return nil, "", *qry.err
	}

	if qry.limit != 1 && qry.typeOf == HybridQry && qry.opts.readGlobalCache {
		keys, cursor, err := qry.GetKeys()
		if err == nil && len(keys) > 0 {
			keys, err = newLoader(qry.coll).Keys(keys...).GetAll(dsts)
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

// Runs the query and writes the result's first entity to the passed destination.
func (qry *Query) GetFirst(dst interface{}) (err error) {
	return qry.Run().GetOne(dst)
}

// Runs the query and returns an Iterator.
func (qry *Query) Run() *Iterator {
	qry.coll.store.ctx.Infof(qry.getLog())
	return &Iterator{qry, qry.qry.Run(qry.coll.store.ctx)}
}

func (qry *Query) log(s string, values ...interface{}) {
	qry.logs = append(qry.logs, fmt.Sprintf(s, values...))
}

func (qry *Query) getLog() string {
	return fmt.Sprintf("running query \"%v\"", strings.Join(qry.logs, " | "))
}
