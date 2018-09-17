package service

import (
	"fmt"

	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/dao"
)

// 域名记录缓存. [domain][type]RR
var rrCache map[string]map[uint16][]*RR

// 域名记录分组缓存. [domain][group]RR
var groupCache map[string]map[string][]*RR

func LoadRRCache() (int, error) {
	sleep := ServConfig.CacheTTL

	rr, err := dao.LoadRRFromMysql()
	if err != nil {
		return sleep, err
	}

	gcache := make(map[string]map[string][]*RR)
	rcache := make(map[string]map[uint16][]*RR)

	for _, r := range rr {
		if int(r.TTL) < sleep {
			sleep = int(r.TTL)
		}

		Logger.Debug(fmt.Sprintf("LoadRRCache: %#v\n", r))

		// gcache
		_, ok := gcache[r.Domain]
		if !ok {
			gcache[r.Domain] = make(map[string][]*RR)
		}
		gcache[r.Domain][r.View] = append(gcache[r.Domain][r.View], r)

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

func GetRRByView(domain, view string) []*RR {
	return groupCache[domain][view]
}
