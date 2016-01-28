package cache

// 每个agent心跳上来的时候立马更新一下数据库是没必要的
// 缓存起来，每隔一个小时写一次DB
// 提供http接口查询机器信息，排查重名机器的时候比较有用

import (
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/hbs/db"
    "github.com/open-falcon/hbs/g"
	"sync"
    "log"
	"time"
)

type SafeAgents struct {
	sync.RWMutex
	M map[string]*model.AgentUpdateInfo
}

var Agents = NewSafeAgents()

func NewSafeAgents() *SafeAgents {
	return &SafeAgents{M: make(map[string]*model.AgentUpdateInfo)}
}

func (this *SafeAgents) Put(req *model.AgentReportRequest) {
	val := &model.AgentUpdateInfo{
		LastUpdate:    time.Now().Unix(),
		ReportRequest: req,
	}
    
    if g.Config().ExternalNodes == "" {
        db.UpdateAgent(val)
    }

	this.Lock()
	defer this.Unlock()
	this.M[req.Hostname] = val
    
    debug := g.Config().Debug
    if debug {
        log.Printf("[DEBUG][CACHE] agent.put : %v", this.M)
    }
}

func (this *SafeAgents) Get(hostname string) (*model.AgentUpdateInfo, bool) {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[hostname]
    debug := g.Config().Debug
    if debug {
        log.Printf("[DEBUG][CACHE] agent.get : hostname:%v, val is %v, exists is %v", hostname,val,exists)
    }
	return val, exists
}

func (this *SafeAgents) Delete(hostname string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, hostname)

    debug := g.Config().Debug
    if debug {
        log.Printf("[DEBUG][CACHE] agent.delete : %v", this.M)
    }
}

func (this *SafeAgents) Keys() []string {
	this.RLock()
	defer this.RUnlock()
	count := len(this.M)
	keys := make([]string, count)
	i := 0
	for hostname := range this.M {
		keys[i] = hostname
		i++
	}

    debug := g.Config().Debug
    if debug {
        log.Printf("[DEBUG][CACHE] agent.keys : %v", keys)
    }
	return keys
}

func DeleteStaleAgents() {
	duration := time.Hour * time.Duration(24)
	for {
		time.Sleep(duration)
		deleteStaleAgents()
	}
}

func deleteStaleAgents() {
	// 一天都没有心跳的Agent，从内存中干掉
	before := time.Now().Unix() - 3600*24
	keys := Agents.Keys()
	count := len(keys)
	if count == 0 {
		return
	}

	for i := 0; i < count; i++ {
		curr, _ := Agents.Get(keys[i])
		if curr.LastUpdate < before {
			Agents.Delete(curr.ReportRequest.Hostname)
		}
	}
}
