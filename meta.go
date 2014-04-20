package hrd

import (
	"fmt"
)

type meta struct {
	stats
	descr string
}

const (
	sourceMemory    = "memory"
	sourceMemcache  = "mcache"
	sourceDatastore = "dstore"
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
			ret = "[" + sourceMemory + "]"
		} else if s.fromDatastore == total {
			ret = "[" + sourceDatastore + "]"
		} else if s.fromGlobalCache == total {
			ret = "[" + sourceMemcache + "]"
		} else {
			ret = "[ "
			if s.fromLocalCache > 0 {
				ret += fmt.Sprintf(sourceMemory+" %d", s.fromLocalCache/total*100.0) + "%% "
			}
			if s.fromDatastore > 0 {
				ret += fmt.Sprintf(sourceDatastore+" %d", s.fromDatastore/total*100.0) + "%% "
			}
			if s.fromGlobalCache > 0 {
				ret += fmt.Sprintf(sourceMemcache+" %d", s.fromGlobalCache/total*100.0) + "%% "
			}
			ret += "]"
		}
	}
	return
}
