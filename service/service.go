package service

import (
	"fmt"

	. "../common"
	"../dao"
	"../geo"
)

// 解析规则缓存. [domain][line/area]Rule
var ruleCache map[string]map[string]*Rule

// rr : [domain][type]RR
func ResolvDomains(clientIP string, rr map[string]map[uint16]*RR) error {
	line, area := geo.FindIP(clientIP)
	Logger.Info("resolv client: " + clientIP + "\t" + line + "\t" + area)

	for domain := range rr {
		rule, ok := ruleCache[domain][line+"/"+area]
		if !ok {
			rule, ok = ruleCache[domain][line+"/*"] //线路默认
			if !ok {
				rule, ok = ruleCache[domain]["*/*"] //域名默认
			}
		}

		Logger.Info("resolv client: " + clientIP + "\t" + line + "\t" + area + ": " + fmt.Sprintf("%#v", rule))

		if rule != nil {
			// 有定义就近解析规则就按规则解析
			r := GroupLB(domain, rule.Group)
			if r == nil {
				Logger.Warn("domain record nil: " + rule.Group)
				continue
			}

			rr[r.Domain][r.Type] = r
			Logger.Debug("resolv one: " + r.Domain)
		} else {
			// 没有定义解析规则
			// @@@　要不要递归权威ＤＮＳ呢？
		}
	}

	Logger.Info(fmt.Sprintf("Resolved: %#v", rr))

	return nil
}

func LoadRuleCache() (int, error) {
	rr, err := dao.CacheRulesFromMysql()
	if err != nil {
		return ServConfig.CacheTTL, err
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

	return ServConfig.CacheTTL, nil
}
