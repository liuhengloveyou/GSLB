package service

import (
	"fmt"

	. "../common"
	"../dao"
)

// 系统中所有接入点分级均衡策略列表.
// rr : 轮询
// wrr: 加权轮询
var lbs []LB = make([]LB, 2)

// 域名记录缓存. [domain][type]RR
var rrCache map[string]map[uint16][]*RR

// 域名记录分组缓存. [domain][group]RR
var groupCache map[string]map[string][]*RR

// 接入分组负载策略配置缓存. ["domain/group"]LB
var groupPolicy map[string]LB

func LoadGroupCache() (int, error) {
	g, err := dao.LoadGroupFromMysql()
	if err != nil {
		return 0, err
	}

	t := make(map[string]LB, len(g))
	for i := 0; i < len(g); i++ {
		if (g[i].Policy) >= len(lbs) {
			Logger.Error(fmt.Sprintf("group policy config ERR: %s %s %d", g[i].Domain, g[i].Name, g[i].Policy))
		}

		t[g[i].Domain+"/"+g[i].Name] = lbs[g[i].Policy]
		Logger.Info(fmt.Sprintf("group one: %v %v %v", g[i].Domain, g[i].Name, lbs[g[i].Policy]))
	}

	return ServConfig.CacheTTL, nil
}

func LoadRRCache() (int, error) {
	sleep := ServConfig.CacheTTL

	rr, err := dao.LoadRRFromMysql()
	if err != nil {
		return sleep, err
	}

	gcache := make(map[string]map[string][]*RR)
	rcache := make(map[string]map[uint16][]*RR)

	for _, r := range rr {
		if int(r.Ttl) < sleep {
			sleep = int(r.Ttl)
		}

		// gcache
		_, ok := gcache[r.Domain]
		if !ok {
			gcache[r.Domain] = make(map[string][]*RR)
		}
		gcache[r.Domain][r.Group] = append(gcache[r.Domain][r.Group], r)

		// rcahce
		_, ok = rcache[r.Domain]
		if !ok {
			rcache[r.Domain] = make(map[uint16][]*RR)
		}
		rcache[r.Domain][r.Type] = append(rcache[r.Domain][r.Type], r)

	}

	rrCache = rcache
	groupCache = gcache
	Logger.Info(fmt.Sprintf("rrCache: %#v\n", rrCache))
	Logger.Info(fmt.Sprintf("groupCache: %#v\n", groupCache))

	return sleep, nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type LB interface {
	Get(domain, group string) *RR
}

type LBRR struct {
	Name string
	i    map[string]int
}

func (p *LBRR) Get(domain, group string) *RR {
	if p.i == nil {
		p.i = make(map[string]int)
		p.i[domain+"/"+group] = -1
	}
	rrs, ok := groupCache[domain][group]
	if !ok {
		return nil // 至少得有默认解析配置
	}

	p.i[domain+"/"+group] = p.i[domain+"/"+group] + 1
	if p.i[domain+"/"+group] >= len(rrs) {
		p.i[domain+"/"+group] = 0
	}

	return rrs[p.i[domain+"/"+group]]
}

type LBWRR struct {
	Name string
}

func (p *LBWRR) Get(domain, group string) *RR {
	rrs, ok := groupCache[domain][group]
	if !ok {
		return nil // 至少得有默认解析配置
	}

	index := -1
	var total int32

	for i := 0; i < len(rrs); i++ {
		rrs[i].CurrentWeight = rrs[i].CurrentWeight + rrs[i].Weight
		total = total + rrs[i].CurrentWeight

		if index == -1 || rrs[index].CurrentWeight < rrs[i].CurrentWeight {
			index = i
		}
	}

	rrs[index].CurrentWeight -= total

	return rrs[index]
}

func GroupLB(domain, group string) *RR {
	if lb, ok := groupPolicy[domain+"/"+group]; ok {
		return lb.Get(domain, group)
	}

	return nil
}

func init() {
	lbs[0] = &LBRR{Name: "rr"}
	lbs[1] = &LBWRR{Name: "wrr"}
}
