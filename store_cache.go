package hrd

import (
	"appengine/memcache"
	"bytes"
	"encoding/gob"
	"github.com/101loops/mcache"
)

// cache manages the caching of a Store.
type cache struct {
	store      *Store
	localCache *mcache.MCache

	// toDelete contains keys to delete after a successful transaction.
	toDelete []*Key

	// toPut contains entities to cache after a successful transaction.
	toPut map[*Key]*doc
}

func newStoreCache(store *Store) *cache {
	if store.tx {
		return &cache{
			store: store,
			toPut: make(map[*Key]*doc),
		}
	}

	return &cache{
		store:      store,
		localCache: mcache.NewMemoryCache(false),
	}
}

func (c *cache) writeTo(dst *cache) {
	var toDelete []*Key
	for _, k := range c.toDelete {
		if _, ok := c.toPut[k]; !ok {
			toDelete = append(toDelete, k)
		}
	}
	dst.delete(toDelete)
	dst.write(c.toPut)
}

func (c *cache) write(toCache map[*Key]*doc) {
	if len(toCache) > 0 {
		if c.store.tx {
			for k, doc := range toCache {
				if k.opts.writeLocalCache || k.opts.writeGlobalCache >= 0 {
					c.toPut[k] = doc
				}
			}
		} else {
			for key, doc := range toCache {
				if key.opts.writeLocalCache {
					c.putMemory(key, doc)
				}
			}
			c.putMemcache(toCache)
		}
	}
}

func (c *cache) delete(keys []*Key) {
	if len(keys) > 0 {
		if c.store.tx {
			for _, k := range keys {
				c.toDelete = append(c.toDelete, k)
			}
		} else {
			i := 0
			memKeys := make([]string, len(keys))
			for _, k := range keys {
				memKey := toMemKey(k)
				c.localCache.Delete(memKey)
				memKeys[i] = memKey
				i++
			}

			memcache.DeleteMulti(c.store.ctx, memKeys) // ignore errors
		}
	}
}

func (c *cache) read(keys []*Key, docs *docs) (dsKeys []*Key, dsDocs []*doc) {

	// #1 populate result from local cache (memory)
	var memIds []int
	var memKeys []string

	for i, key := range keys {
		mKey := toMemKey(key)

		// lookup key in memory
		if key.opts.readLocalCache && !c.store.tx {
			if src, ok := c.getMemory(key); ok && src != nil {
				key.source = sourceMemory
				docs.set(i, src)
				continue
			}
		}

		// remember missing keys for following lookup
		memIds = append(memIds, i)
		memKeys = append(memKeys, mKey)
	}

	// #2 populate result from global cache (memcache)
	memVals, err := memcache.GetMulti(c.store.ctx, memKeys)
	if err != nil {
		c.store.logErr(err)
	}

	for i, mKey := range memKeys {
		dstID := memIds[i]
		key := keys[dstID]
		doc := docs.get(dstID)

		// lookup key in global cache result
		if key.opts.readGlobalCache && !c.store.tx {
			if item, ok := memVals[mKey]; ok && item.Value != nil {
				if err := fromGob(doc, item.Value); err == nil {
					key.source = sourceMemcache
					if key.opts.writeLocalCache {
						c.putMemory(key, doc) // copy to local cache as well
					}
					continue
				} else {
					c.store.logErr(err)
				}
			}
		}

		// remember missing keys for following lookup
		dsDocs = append(dsDocs, doc)
		dsKeys = append(dsKeys, key)
	}

	return
}

func (c *cache) putMemcache(toPut map[*Key]*doc) {
	var items []*memcache.Item
	for key, doc := range toPut {
		expire := key.opts.writeGlobalCache
		if expire == -1 {
			continue
		}

		gob, err := toGob(doc)
		if err == nil {
			items = append(items, &memcache.Item{
				Key:        toMemKey(key),
				Value:      gob,
				Expiration: key.opts.writeGlobalCache,
			})
		} else {
			c.store.logErr(err)
		}
	}
	if len(items) > 0 {
		err := memcache.SetMulti(c.store.ctx, items)
		if err != nil {
			c.store.logErr(err)
		}
	}
}

func (c *cache) getMemory(key *Key) (interface{}, bool) {
	return c.localCache.Get(toMemKey(key))
}

func (c *cache) putMemory(key *Key, doc *doc) {
	c.localCache.PutP(toMemKey(key), doc.get())
}

func toGob(doc *doc) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	src := doc.get()
	gob.Register(src)

	if err := enc.Encode(src); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func fromGob(doc *doc, b []byte) error {
	buf := bytes.NewBuffer(b)

	src := doc.get()
	gob.Register(src)

	dec := gob.NewDecoder(buf)
	if err := dec.Decode(src); err != nil {
		return err
	}
	doc.set(src)

	return nil
}
