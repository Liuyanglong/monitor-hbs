package cache

import (
	"github.com/open-falcon/hbs/db"
	"github.com/open-falcon/hbs/enc"
	"github.com/open-falcon/hbs/g"
	"log"
	"sync"
)

// 一个机器可能在多个group下，做一个map缓存hostid与groupid的对应关系
type SafeHostGroupsMap struct {
	sync.RWMutex
	M map[int][]int
}

var HostGroupsMap = &SafeHostGroupsMap{M: make(map[int][]int)}

func (this *SafeHostGroupsMap) GetGroupIds(hid int) ([]int, bool) {
	this.RLock()
	defer this.RUnlock()
	gids, exists := this.M[hid]
	debug := g.Config().Debug
	if debug {
		log.Printf("[DEBUG][CACHE] group.GetGroupIds :hid is %v, gids is %v, exists is %v", hid, gids, exists)
	}
	return gids, exists
}

func (this *SafeHostGroupsMap) Init() {
	var hostGroupMap map[int][]int

	if g.Config().ExternalNodes == "" {
		m, err := db.QueryHostGroups()
		if err != nil {
			return
		}
		hostGroupMap = m
	} else {
		m, err := enc.QueryHostGroups()
		if err != nil {
			return
		}
		hostGroupMap = m
	}

	this.Lock()
	defer this.Unlock()
	this.M = hostGroupMap
	debug := g.Config().Debug
	if debug {
		log.Printf("[DEBUG][CACHE] group.init : %v", this.M)
	}
}
