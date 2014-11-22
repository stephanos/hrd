package types

import (
	ae "appengine"
)

// Kind represents a entity category in the datastore.
type Kind struct {
	Context ae.Context
	Name    string
}

// NewKind creates a new kind.
func NewKind(ctx ae.Context, name string) *Kind {
	return &Kind{ctx, name}
}

//// NewNumKey returns a key for the passed kind and numeric ID.
//// It can also receive an optional parent key.
//func (s *Store) NewNumKey(kind string, id int64, parent ...*Key) *Key {
//	var parentKey *ds.Key
//	if len(parent) > 0 {
//		parentKey = parent[0].Key.Key
//	}
//	return newKey(internal.NewKey(ds.NewKey(s.ctx, kind, "", id, parentKey)))
//}
//
//// NewNumKeys returns a sequence of key for the passed kind and
//// sequence of numeric ID.
//func (s *Store) NewNumKeys(kind string, ids ...int64) []*Key {
//	keys := make([]*Key, len(ids))
//	for i, id := range ids {
//		keys[i] = s.NewNumKey(kind, id)
//	}
//	return keys
//}
//
//// NewTextKey returns a key for the passed kind and string ID.
//// It can also receive an optional parent key.
//func (s *Store) NewTextKey(kind string, id string, parent ...*Key) *Key {
//	var parentKey *ds.Key
//	if len(parent) > 0 {
//		parentKey = parent[0].Key.Key
//	}
//	return newKey(internal.NewKey(ds.NewKey(s.ctx, kind, id, 0, parentKey)))
//}
//
//// NewTextKeys returns a sequence of keys for the passed kind and
//// sequence of string ID.
//func (s *Store) NewTextKeys(kind string, ids ...string) []*Key {
//	keys := make([]*Key, len(ids))
//	for i, id := range ids {
//		keys[i] = s.NewTextKey(kind, id)
//	}
//	return keys
//}
