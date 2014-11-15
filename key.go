package hrd

import (
	"github.com/101loops/hrd/internal"
)

// Key represents the datastore key for an
// It also contains meta data about said
type Key struct {
	*internal.Key
}

func newKey(k *internal.Key) *Key {
	return &Key{Key: k}
}

func newKeys(keys []*internal.Key) []*Key {
	ret := make([]*Key, len(keys))
	for i, k := range keys {
		ret[i] = newKey(k)
	}
	return ret
}

func toInternalKeys(keys []*Key) []*internal.Key {
	ret := make([]*internal.Key, len(keys))
	for i, k := range keys {
		ret[i] = k.Key
	}
	return ret
}
