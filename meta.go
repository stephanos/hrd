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

func (self *meta) string() (s string) {
	s = self.descr
	if s != "" {
		s += " "
	}
	s += self.stats.string()
	return
}

type stats struct {
	fromLocalCache  int
	fromGlobalCache int
	fromDatastore   int
}

func (self *stats) string() (s string) {
	total := int(self.fromLocalCache + self.fromGlobalCache + self.fromDatastore)
	if total > 0 {
		if self.fromLocalCache == total {
			s = "[" + SOURCE_MEMORY + "]"
		} else if self.fromDatastore == total {
			s = "[" + SOURCE_DATASTORE + "]"
		} else if self.fromGlobalCache == total {
			s = "[" + SOURCE_MEMCACHE + "]"
		} else {
			s = "[ "
			if self.fromLocalCache > 0 {
				s += fmt.Sprintf(SOURCE_MEMORY+" %d", self.fromLocalCache/total*100.0) + "%% "
			}
			if self.fromDatastore > 0 {
				s += fmt.Sprintf(SOURCE_DATASTORE+" %d", self.fromDatastore/total*100.0) + "%% "
			}
			if self.fromGlobalCache > 0 {
				s += fmt.Sprintf(SOURCE_MEMCACHE+" %d", self.fromGlobalCache/total*100.0) + "%% "
			}
			s += "]"
		}
	}
	return
}
