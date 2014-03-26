package hrd

import (
	"fmt"
)

type meta struct {
	stats
	descr string
}

const (
	SOURCE_MEMORY    = "memory"
	SOURCE_MEMCACHE  = "mcache"
	SOURCE_DATASTORE = "dstore"
)

func (m *meta) string() (ret string) {
	ret = m.descr
	if ret != "" {
		ret += " "
	}
	ret += m.stats.string()
	return
}

type stats struct {
	fromLocalCache  int
	fromGlobalCache int
	fromDatastore   int
}

func (s *stats) string() (ret string) {
	total := int(s.fromLocalCache + s.fromGlobalCache + s.fromDatastore)
	if total > 0 {
		if s.fromLocalCache == total {
			ret = "[" + SOURCE_MEMORY + "]"
		} else if s.fromDatastore == total {
			ret = "[" + SOURCE_DATASTORE + "]"
		} else if s.fromGlobalCache == total {
			ret = "[" + SOURCE_MEMCACHE + "]"
		} else {
			ret = "[ "
			if s.fromLocalCache > 0 {
				ret += fmt.Sprintf(SOURCE_MEMORY+" %d", s.fromLocalCache/total*100.0) + "%% "
			}
			if s.fromDatastore > 0 {
				ret += fmt.Sprintf(SOURCE_DATASTORE+" %d", s.fromDatastore/total*100.0) + "%% "
			}
			if s.fromGlobalCache > 0 {
				ret += fmt.Sprintf(SOURCE_MEMCACHE+" %d", s.fromGlobalCache/total*100.0) + "%% "
			}
			ret += "]"
		}
	}
	return
}
