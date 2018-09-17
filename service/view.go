package service

import (
	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/dao"
)

// 解析规则缓存. [domain][line/area]Rule
/*
var viewCache map[string]map[string]*Rule

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

	viewCache = cache
	Logger.Info(fmt.Sprintf("ruleCache: %#v\n", viewCache))

	return ServConfig.CacheTTL, nil
}
*/

func GetView(line, area string) (view *View, e error) {
	return dao.SelectViewFromMysql(line, area)
}
