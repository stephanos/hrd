package types

import (
	ae "appengine"
	ds "appengine/datastore"
)

// Order is a sort order on query results.
type Order struct {
	FieldName  string
	Descending bool
}

// Filter is a conditional filter on query results.
type Filter struct {
	Filter string
	Value  interface{}
}

// QueryType describes the way a query is run and what it returns.
type QueryType int

const (
	// FullQuery is a regular query.
	FullQuery QueryType = 1 + iota

	// ProjectQuery is a query that only yields selected fields.
	ProjectQuery

	// KeysOnlyQuery is a query that only returns the result's keys.
	KeysOnlyQuery
)

// Query represents a datastore query.
type Query struct {
	kind string

	Ancestor   *Key
	Filter     []Filter
	Order      []Order
	Projection []string

	TypeOf   QueryType
	Distinct bool
	Eventual bool
	Limit    int
	Offset   int
	Start    string
	End      string
}

// NewQuery creates a new, empty query.
func NewQuery(kind string) (ret *Query) {
	return &Query{
		Filter: make([]Filter, 0),
		Order:  make([]Order, 0),
		TypeOf: FullQuery,
		kind:   kind,
		Limit:  -1,
	}
}

// Clone creates a deep copy.
func (q *Query) Clone() *Query {
	ret := *q

	ret.Order = make([]Order, len(q.Order))
	copy(ret.Order, q.Order)

	ret.Filter = make([]Filter, len(q.Filter))
	copy(ret.Filter, q.Filter)

	ret.Projection = make([]string, len(q.Projection))
	copy(ret.Projection, q.Projection)

	return &ret
}

// ToDSQuery converts the query to a datastore Query.
func (q *Query) ToDSQuery(ctx ae.Context) *ds.Query {
	dsQry := ds.NewQuery(q.kind).Limit(q.Limit)
	for _, f := range q.Filter {
		dsQry = dsQry.Filter(f.Filter, f.Value)
	}
	for _, o := range q.Order {
		order := o.FieldName
		if o.Descending {
			order = "-" + order
		}
		dsQry = dsQry.Order(order)
	}
	if len(q.Projection) > 0 {
		dsQry = dsQry.Project(q.Projection...)
	}
	if q.Ancestor != nil {
		dsQry = dsQry.Ancestor(q.Ancestor.ToDSKey(ctx))
	}
	if q.Start != "" {
		cursor, _ := ds.DecodeCursor(q.Start)
		dsQry = dsQry.Start(cursor)
	}
	if q.End != "" {
		cursor, _ := ds.DecodeCursor(q.End)
		dsQry = dsQry.End(cursor)
	}
	if q.Offset != 0 {
		dsQry = dsQry.Offset(q.Offset)
	}
	if q.Eventual {
		dsQry = dsQry.EventualConsistency()
	}
	if q.Distinct {
		dsQry = dsQry.Distinct()
	}
	if q.TypeOf == KeysOnlyQuery {
		dsQry = dsQry.KeysOnly()
	}
	return dsQry
}
