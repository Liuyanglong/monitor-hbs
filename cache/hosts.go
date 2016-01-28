package cache

import (
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/hbs/db"
	"github.com/open-falcon/hbs/enc"
	"github.com/open-falcon/hbs/g"
	"log"
	"sync"
)

// 每次心跳的时候agent把hostname汇报上来，经常要知道这个机器的hostid，把此信息缓存
// key: hostname value: hostid
type SafeHostMap struct {
	sync.RWMutex
	M map[string]int
}

var HostMap = &SafeHostMap{M: make(map[string]int)}

func (this *SafeHostMap) GetID(hostname string) (int, bool) {
	this.RLock()
	defer this.RUnlock()
	id, exists := this.M[hostname]

	debug := g.Config().Debug
	if debug {
		log.Printf("[DEBUG][CACHE] host.getid : hostname is %v,id is %v, exists is %v", hostname, id, exists)
	}
	return id, exists
}

func (this *SafeHostMap) Init() {
	var hostMap map[string]int

	if g.Config().ExternalNodes == "" {
		m, err := db.QueryHosts()
		if err != nil {
			return
		}
		hostMap = m
	} else {
		m, err := enc.QueryHosts()
		if err != nil {
			return
		}
		hostMap = m
	}

	this.Lock()
	defer this.Unlock()
	this.M = hostMap
	debug := g.Config().Debug
	if debug {
		log.Printf("[DEBUG][CACHE] host.init : %v", this.M)
	}
}

type SafeMonitoredHosts struct {
	sync.RWMutex
	M map[int]*model.Host
}

var MonitoredHosts = &SafeMonitoredHosts{M: make(map[int]*model.Host)}

func (this *SafeMonitoredHosts) Get() map[int]*model.Host {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

//暂时没有实现基于enc的此方法
func (this *SafeMonitoredHosts) Init() {
	var hostMap map[int]*model.Host

	if g.Config().ExternalNodes == "" {
		m, err := db.QueryMonitoredHosts()
		if err != nil {
			return
		}
		hostMap = m
	} else {
		m, err := enc.QueryMonitoredHosts()
		if err != nil {
			return
		}
		hostMap = m
	}

	this.Lock()
	defer this.Unlock()
	this.M = hostMap

	debug := g.Config().Debug
	if debug {
		log.Printf("[DEBUG][CACHE] SafeMonitoredHosts.init : %v", this.M)
	}
}
