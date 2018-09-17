package service

import (
	"fmt"

	. "../common"
	"../dao"
)

// 域名记录缓存. [domain][type]RR
var rrCache map[string]map[string][]RR

// 域名记录分组缓存. [domain][group]RR
var groupCache map[string]map[string][]RR

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

	gcache := make(map[string]map[string][]RR)
	rcache := make(map[string]map[string][]RR)

	for _, r := range rr {
		if int(r.TTL) < sleep {
			sleep = int(r.TTL)
		}

		// gcache
		_, ok := gcache[r.Domain]
		if !ok {
			gcache[r.Domain] = make(map[string][]RR)
		}
		gcache[r.Domain][r.View] = append(gcache[r.Domain][r.View], r)

		// rcahce
		_, ok = rcache[r.Domain]
		if !ok {
			rcache[r.Domain] = make(map[string][]RR)
		}
		rcache[r.Domain][r.Type] = append(rcache[r.Domain][r.Type], r)

	}

	rrCache = rcache
	groupCache = gcache
	Logger.Info(fmt.Sprintf("rrCache: %#v\n", rrCache))
	Logger.Info(fmt.Sprintf("groupCache: %#v\n", groupCache))

	return sleep, nil
}

func GetRRByView(domain, view string) []RR {
	return groupCache[domain][view]
}
