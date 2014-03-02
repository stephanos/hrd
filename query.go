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

func newQuery(coll *Collection) *Query {
	return &Query{
		coll:   coll,
		limit:  -1,
		typeOf: HybridQry,
		opts:   defaultOperationOpts(),
		qry:    datastore.NewQuery(coll.name),
		logs:   []string{"KIND " + coll.name},
	}
}

func (self *Query) clone() *Query {
	q := *self
	q.opts = self.opts.clone()
	if len(self.logs) > 0 {
		q.logs = make([]string, len(self.logs))
		copy(q.logs, self.logs)
	}
	return &q
}

func (self *Query) Flags(flags ...Flag) *Query {
	q := self.clone()
	q.opts = q.opts.Flags(flags...)
	return q
}

func (self *Query) NoHybrid() *Query {
	return self.Hybrid(false)
}

func (self *Query) Hybrid(enabled bool) *Query {
	q := self.clone()
	if enabled {
		if q.typeOf == FullQry {
			q.typeOf = HybridQry
		}
	} else {
		if q.typeOf == HybridQry {
			q.typeOf = FullQry
		}
	}
	return q
}

func (self *Query) Limit(limit int) *Query {
	q := self.clone()
	if limit > 0 {
		q.log("LIMIT %v", limit)
	} else {
		limit = -1
		q.log("NO LIMIT")
	}
	q.qry = q.qry.Limit(limit)
	q.limit = limit
	return q
}

func (self *Query) NoLimit() *Query {
	return self.Limit(-1)
}

func (self *Query) Ancestor(k *Key) *Query {
	q := self.clone()
	q.log("ANCESTOR '%v'", k.IdString())
	q.qry = q.qry.Ancestor(k.Key)
	return q
}

func (self *Query) Project(s ...string) *Query {
	q := self.clone()
	q.log("PROJECT '%v'", strings.Join(s, "', '"))
	q.qry = q.qry.Project(s...)
	q.typeOf = ProjQry
	return q
}

func (self *Query) End(c string) *Query {
	q := self.clone()
	if c != "" {
		if cursor, err := datastore.DecodeCursor(c); err == nil {
			q.log("END CURSOR")
			q.qry = q.qry.End(cursor)
		} else {
			err = fmt.Errorf("invalid end cursor (%v)", err)
			q.err = &err
		}
	}
	return q
}

func (self *Query) Start(c string) *Query {
	q := self.clone()
	if c != "" {
		if cursor, err := datastore.DecodeCursor(c); err == nil {
			q.log("START CURSOR")
			q.qry = q.qry.Start(cursor)
		} else {
			err = fmt.Errorf("invalid start cursor (%v)", err)
			q.err = &err
		}
	}
	return q
}

func (self *Query) Offset(off int) *Query {
	q := self.clone()
	q.log("OFFSET %v", off)
	q.qry = q.qry.Offset(off)
	return q
}

func (self *Query) OrderAsc(s string) *Query {
	q := self.clone()
	q.log("ORDER ASC %v", s)
	q.qry = q.qry.Order(s)
	return q
}

func (self *Query) OrderDesc(s string) *Query {
	q := self.clone()
	q.log("ORDER DESC %v", s)
	q.qry = q.qry.Order("-" + s)
	return q
}

func (self *Query) Filter(qry string, val interface{}) *Query {
	q := self.clone()
	q.log("FILTER '%v %v'", qry, val)
	q.qry = q.qry.Filter(qry, val)
	return q
}

// ==== CACHE

func (self *Query) NoCache() *Query {
	return self.NoLocalCache().NoGlobalCache()
}

func (self *Query) NoLocalCache() *Query {
	return self.NoLocalCacheWrite().NoLocalCacheRead()
}

func (self *Query) NoGlobalCache() *Query {
	return self.NoGlobalCacheWrite().NoGlobalCacheRead()
}

func (self *Query) CacheExpire(exp time.Duration) *Query {
	q := self.clone()
	q.opts = q.opts.CacheExpire(exp)
	return q
}

func (self *Query) NoCacheRead() *Query {
	return self.NoGlobalCacheRead().NoLocalCacheRead()
}

func (self *Query) NoLocalCacheRead() *Query {
	q := self.clone()
	q.opts = q.opts.NoLocalCacheRead()
	return q
}

func (self *Query) NoGlobalCacheRead() *Query {
	q := self.clone()
	q.opts = q.opts.NoGlobalCacheRead()
	return q
}

func (self *Query) NoCacheWrite() *Query {
	return self.NoGlobalCacheWrite().NoLocalCacheWrite()
}

func (self *Query) NoLocalCacheWrite() *Query {
	q := self.clone()
	q.opts = q.opts.NoLocalCacheWrite()
	return q
}

func (self *Query) NoGlobalCacheWrite() *Query {
	q := self.clone()
	q.opts = q.opts.NoGlobalCacheWrite()
	return q
}

// ==== EXECUTE

func (self *Query) GetCount() (int, error) {
	self.log("COUNT")
	self.coll.store.ctx.Infof(self.getLog())

	if self.err != nil {
		return 0, *self.err
	}
	return self.qry.Count(self.coll.store.ctx)
}

// Runs the query as keys-only: No entities are retrieved, just their keys.
func (self *Query) GetKeys() ([]*Key, string, error) {
	q := self.clone()
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
func (self *Query) GetAll(dsts interface{}) ([]*Key, string, error) {
	if self.err != nil {
		return nil, "", *self.err
	}

	if self.limit != 1 && self.typeOf == HybridQry && self.opts.readGlobalCache {
		keys, cursor, err := self.GetKeys()
		if err == nil && len(keys) > 0 {
			keys, err = newLoader(self.coll).Keys(keys...).GetAll(dsts)
		}
		return keys, cursor, err
	}

	it := self.Run()
	keys, err := it.GetAll(dsts)
	if err != nil {
		return nil, "", err
	}

	cursor, err := it.Cursor()
	return keys, cursor, err
}

// Runs the query and writes the result's first entity to the passed destination.
func (self *Query) GetFirst(dst interface{}) (err error) {
	return self.Run().GetOne(dst)
}

// Runs the query and returns an Iterator.
func (self *Query) Run() *Iterator {
	self.coll.store.ctx.Infof(self.getLog())
	return &Iterator{self, self.qry.Run(self.coll.store.ctx)}
}

func (self *Query) log(s string, values ...interface{}) {
	self.logs = append(self.logs, fmt.Sprintf(s, values...))
}

func (self *Query) getLog() string {
	return fmt.Sprintf("running query \"%v\"", strings.Join(self.logs, " | "))
}
