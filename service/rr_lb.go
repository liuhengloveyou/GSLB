package service

import (
	"sort"

	"../common"
	"../geo"
)

// 系统中所有接入点分级均衡策略列表.
// rr : 轮询
// wrr: 加权轮询
var lbs []LB = make([]LB, 3)

// 接入View负载策略配置缓存. ["domain/view"]LB
var groupPolicy map[string]LB

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type LB interface {
	Get(domain, view string, client *geo.IpRecord, rrs []common.RR) []common.RR
}

// 轮询
type LBRR struct {
	Name string
	i    map[string]int
}

func (p *LBRR) Get(domain, view string, client *geo.IpRecord, rrs []common.RR) []common.RR {
	if p.i == nil {
		p.i = make(map[string]int)
		p.i[domain+"/"+view] = -1
	}

	p.i[domain+"/"+view] = p.i[domain+"/"+view] + 1

	if p.i[domain+"/"+view] >= len(rrs) {
		p.i[domain+"/"+view] = 0
	}

	idx := p.i[domain+"/"+view]
	rrsr := make([]common.RR, len(rrs))
	for i := 0; i < len(rrs); i++ {
		if idx >= len(rrs) {
			idx = 0
		}
		rrsr[i] = rrs[idx]
	}

	return rrsr
}

// 加权轮询
type LBWRR struct {
	Name string
}

func (p *LBWRR) Get(domain, view string, client *geo.IpRecord, rrs []common.RR) []common.RR {
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

	rrsr := make([]common.RR, len(rrs))
	for i := 0; i < len(rrs); i++ {
		rrsr[i] = rrs[index]
		if index >= len(rrs) {
			index = 0
		}
	}

	return rrsr
}

// 就近接入
type ByDistance []common.RR

func (p ByDistance) Len() int           { return len(p) }
func (p ByDistance) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByDistance) Less(i, j int) bool { return p[i].Distance < p[j].Distance }

type LBNear struct {
	Name string
}

func (p *LBNear) Get(domain, view string, client *geo.IpRecord, rrs []common.RR) []common.RR {
	var rrsr = make([]common.RR, 0)

	for i := 0; i < len(rrs); i++ {
		if rrs[i].Type == "A" {
			ip, _ := geo.FindIP(rrs[i].Record)
			rrs[i].Distance = geo.LatitudeLongitudeDistance(client.Latitude, client.Longitude, ip.Latitude, ip.Longitude)
			rrsr = append(rrsr, rrs[i])
		}
	}

	if len(rrsr) > 0 {
		sort.Sort(ByDistance(rrsr))
	} else {
		rrsr = rrs
	}

	return rrsr
}

func GroupLB(domain, view string, client *geo.IpRecord, rrs []common.RR) []common.RR {
	if lb, ok := groupPolicy[domain+"/"+view]; ok {
		return lb.Get(domain, view, client, rrs)
	} else {
		return lbs[0].Get(domain, view, client, rrs) // 默认
	}

	return nil
}

func init() {
	lbs[0] = &LBNear{Name: "near"}
	lbs[1] = &LBRR{Name: "rr"}
	lbs[2] = &LBWRR{Name: "wrr"}
}
