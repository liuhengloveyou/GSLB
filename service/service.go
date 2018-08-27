package service

import (
	"fmt"
	"time"

	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/dao"
	"github.com/liuhengloveyou/GSLB/geo"
)

// 域名记录缓存
// [domain][type]RR
var rrCache map[string]map[uint16][]*RR

// 域名记录分组缓存
// [domain][group]RR
var groupCache map[string]map[string][]*RR

// 解析规则缓存
// [domain][line/area]Rule
var ruleCache map[string]map[string]*Rule

func groupLB(domain, group string) *RR {
	rrs, ok := groupCache[domain][group]
	if !ok {
		return nil
	}

	p := int64(time.Now().Unix()) % int64(len(rrs))

	return rrs[p]
}

// rr : [domain][type]RR
func ResolvDomains(clientIP string, rr map[string]map[uint16]*RR) error {
	// 所有域名
	d := make([]string, len(rr))
	i := 0
	for k := range rr {
		d[i] = k
		i++
	}

	line, area := geo.FindIP(clientIP)
	Logger.Info("resolv client: " + clientIP + "\t" + line + "\t" + area)

	for domain := range rr {
		rule := ruleCache[domain][line+"/"+area]
		Logger.Info("resolv client: " + clientIP + "\t" + line + "\t" + area + ": " + fmt.Sprintf("%#v", rule))

		if rule != nil {
			// 有定义就近解析规则就按规则解析
			r := groupLB(domain, rule.Group)
			if r == nil {
				Logger.Warn("domain record nil: " + rule.Group)
				continue
			}

			rr[r.Domain][r.Type] = r
			Logger.Debug("resolv one: " + r.Domain)
		} else {
			// 没有定义解析规则
			// @@@
		}
	}

	Logger.Info(fmt.Sprintf("Resolved: %#v", rr))

	return nil
}

// SELECT rule.domain, zone.line, zone.area,rule.group FROM rule join zone on rule.zone = zone.zone;
func LoadRuleCache() {
	defer Logger.Sync()

	for {
		rr, err := dao.CacheRulesFromMysql("0")
		if err != nil {
			Logger.Error("LoadRuleCache ERR: " + err.Error())
			time.Sleep(time.Second * time.Duration(ServConfig.CacheTTL))
		}

		cache := make(map[string]map[string]*Rule)
		for _, r := range rr {
			_, ok := cache[r.Domain]
			if !ok {
				cache[r.Domain] = make(map[string]*Rule)
			}
			cache[r.Domain][r.Line+"/"+r.Area] = r
		}

		ruleCache = cache
		Logger.Info(fmt.Sprintf("ruleCache: %#v\n", ruleCache))

		time.Sleep(time.Second * time.Duration(ServConfig.CacheTTL))
	}

}

func LoadRRCache() {
	now := "0"

	defer Logger.Sync()

	for {
		sleep := ServConfig.CacheTTL

		rr, err := dao.LoadRRFromMysql(now)
		if err != nil {
			Logger.Error("LoadRRFromMysql ERR: " + err.Error())
			time.Sleep(time.Second * time.Duration(ServConfig.CacheTTL))
		}

		gcache := make(map[string]map[string][]*RR)
		rcache := make(map[string]map[uint16][]*RR)

		for _, r := range rr {
			if int64(r.Ttl) < sleep {
				sleep = int64(r.Ttl)
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

		//now = time.Now().Format("2006-01-02 15:04:05")
		time.Sleep(time.Second * time.Duration(sleep))
	}
}
