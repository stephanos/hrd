package hrd

import (
	"fmt"
	. "github.com/101loops/bdd"
	"github.com/101loops/hrd/internal"
	"github.com/101loops/hrd/internal/types"

	ae "appengine"
)

var _ = Describe("Query", func() {

	var (
		query *Query
	)

	BeforeEach(func() {
		query = myKind.Query(ctx)
		dsCount = func(_ ae.Context, _ *types.Query) (int, error) {
			panic("unexpected call")
		}
		dsIterate = func(_ *types.Iterator, _ interface{}, _ bool) ([]*types.Key, error) {
			panic("unexpected call")
		}
		dsGet = func(_ *types.Kind, _ []*types.Key, _ interface{}, _ bool, _ bool) ([]*types.Key, error) {
			panic("unexpected call")
		}
	})

	AfterEach(func() {
		dsGet = internal.Get
		dsCount = internal.Count
		dsIterate = internal.Iterate
	})

	Context("building the query", func() {

		It("limit", func() {
			Check(query.inner.Limit, EqualsNum, -1)

			query = query.Limit(10)
			Check(query.inner.Limit, EqualsNum, 10)

			query = query.NoLimit()
			Check(query.inner.Limit, EqualsNum, -1)

			query = query.Limit(0)
			Check(query.inner.Limit, EqualsNum, -1)

			query = query.Limit(-1)
			Check(query.inner.Limit, EqualsNum, -1)

			query = query.Limit(-10)
			Check(query.inner.Limit, EqualsNum, -1)
		})

		It("project", func() {
			Check(query.inner.Projection, IsEmpty)
			Check(query.inner.TypeOf, Equals, types.FullQuery)

			projQuery := query.Project("a", "b", "c")
			Check(projQuery.inner.TypeOf, Equals, types.ProjectQuery)
			Check(projQuery.inner.Projection, Equals, []string{"a", "b", "c"})
		})

		It("start", func() {
			Check(query.inner.Start, IsZero)

			query = query.Start("start-cursor")
			Check(query.inner.Start, Equals, "start-cursor")
		})

		It("end", func() {
			Check(query.inner.End, IsZero)

			query = query.End("end-cursor")
			Check(query.inner.End, Equals, "end-cursor")
		})

		It("offset", func() {
			Check(query.inner.Offset, IsZero)

			query = query.Offset(10)
			Check(query.inner.Offset, EqualsNum, 10)
		})

		It("ancestor", func() {
			Check(query.inner.Ancestor, IsNil)

			key := newNumKey(myKind, 42, nil)
			query = query.Ancestor(key)
			Check(query.inner.Ancestor, Equals, key.inner)
		})

		It("distinct", func() {
			Check(query.inner.Distinct, IsFalse)

			query = query.Distinct()
			Check(query.inner.Distinct, IsTrue)
		})

		It("eventual consistency", func() {
			Check(query.inner.Eventual, IsFalse)

			query = query.EventualConsistency()
			Check(query.inner.Eventual, IsTrue)
		})

		It("order", func() {
			Check(query.inner.Order, IsEmpty)

			query = query.OrderAsc("age")
			Check(query.inner.Order, HasLen, 1)
			Check(query.inner.Order, Contains, types.Order{FieldName: "age", Descending: false})

			query = query.OrderDesc("name")
			Check(query.inner.Order, HasLen, 2)
			Check(query.inner.Order, Contains, types.Order{FieldName: "name", Descending: true})
		})

		It("filter", func() {
			Check(query.inner.Filter, IsEmpty)

			query = query.Filter("age >", 18)
			Check(query.inner.Filter, HasLen, 1)
			Check(query.inner.Filter, Contains, types.Filter{Filter: "age >", Value: 18})

			query = query.Filter("count <", 1000)
			Check(query.inner.Filter, HasLen, 2)
			Check(query.inner.Filter, Contains, types.Filter{Filter: "count <", Value: 1000})
		})
	})

	Context("executing the query", func() {

		var (
			retKeys = []*types.Key{
				types.NewKey("my-kind", "", 1, nil), types.NewKey("my-kind", "", 2, nil),
			}
		)

		It("should return an iterator", func() {
			it := query.Run()
			Check(it, NotNil)
		})

		Context("count", func() {

			It("should return the result's size", func() {
				dsCount = func(_ ae.Context, _ *types.Query) (int, error) {
					return 42, nil
				}
				c, err := query.GetCount()
				Check(c, EqualsNum, 42)
				Check(err, IsNil)
			})

			It("should return an error when the operation fails", func() {
				dsCount = func(_ ae.Context, _ *types.Query) (int, error) {
					return 0, fmt.Errorf("an error")
				}
				c, err := query.GetCount()
				Check(c, EqualsNum, 0)
				Check(err, HasOccurred)
			})
		})

		Context("keys", func() {

			It("should return the result's keys", func() {
				dsIterate = func(_ *types.Iterator, _ interface{}, multi bool) ([]*types.Key, error) {
					Check(multi, IsTrue)
					return retKeys, nil
				}
				keys, _, err := query.GetKeys()
				Check(err, IsNil)
				Check(keys, Equals, importKeys(retKeys))
			})

			It("should return an error when the operation fails", func() {
				dsIterate = func(_ *types.Iterator, _ interface{}, _ bool) ([]*types.Key, error) {
					return nil, fmt.Errorf("an error")
				}

				_, _, err := query.GetKeys()
				Check(err, HasOccurred)
			})
		})

		Context("one entity", func() {

			var entity MyModel

			It("should return the first result", func() {
				dsIterate = func(_ *types.Iterator, _ interface{}, multi bool) ([]*types.Key, error) {
					Check(multi, IsFalse)
					return retKeys[0:1], nil
				}
				key, err := query.GetFirst(&entity)
				Check(err, IsNil)
				Check(key, Equals, importKey(retKeys[0]))
			})

			It("should return nil when there is no result", func() {
				dsIterate = func(_ *types.Iterator, _ interface{}, _ bool) ([]*types.Key, error) {
					return []*types.Key{}, nil
				}
				key, err := query.GetFirst(&entity)
				Check(err, IsNil)
				Check(key, IsNil)
			})
		})

		Context("all entities", func() {

			var entities []*MyModel

			It("should use hybrid query by default", func() {
				dsIterate = func(_ *types.Iterator, _ interface{}, multi bool) ([]*types.Key, error) {
					Check(multi, IsTrue)
					return retKeys, nil
				}

				dsGet = func(kind *types.Kind, keys []*types.Key, _ interface{}, useGlobalCache bool, multi bool) ([]*types.Key, error) {
					Check(kind.Name, Equals, myKind.name)
					Check(keys, Equals, retKeys)
					Check(useGlobalCache, IsTrue)
					Check(multi, IsTrue)
					return retKeys, nil
				}

				keys, _, err := query.GetAll(&entities)
				Check(err, IsNil)
				Check(keys, Equals, importKeys(retKeys))
			})

			It("should run the iterator otherwise", func() {
				fetchWithIterator := func(q *Query) {
					dsIterate = func(_ *types.Iterator, _ interface{}, multi bool) ([]*types.Key, error) {
						Check(multi, IsTrue)
						return retKeys, nil
					}

					keys, _, err := q.GetAll(&entities)
					Check(err, IsNil)
					Check(keys, Equals, importKeys(retKeys))
				}

				fetchWithIterator(query.Limit(1))
				fetchWithIterator(query.Project("a"))
				fetchWithIterator(query.NoGlobalCache())
			})

			It("should return an error when the query fails", func() {
				dsIterate = func(_ *types.Iterator, _ interface{}, _ bool) ([]*types.Key, error) {
					return nil, fmt.Errorf("an error")
				}

				_, _, err := query.NoGlobalCache().GetAll(&entities)
				Check(err, HasOccurred)
			})
		})
	})

})
