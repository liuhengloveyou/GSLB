package service

import (
	"fmt"

	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/dao"
)

// VIEW缓存. [domain][line/area]Rule
var viewCache map[string]map[string]*View

func LoadRuleCache() (int, error) {
	rr, err := dao.LoadViewFromMysql()
	if err != nil {
		return ServConfig.CacheTTL, err
	}

	cache := make(map[string]map[string]*View)
	for _, r := range rr {
		_, ok := cache[r.Domain]
		if !ok {
			cache[r.Domain] = make(map[string]*View)
		}
		cache[r.Domain][r.Line+"/"+r.Area] = r
	}

	viewCache = cache
	Logger.Info(fmt.Sprintf("ruleCache: %#v\n", viewCache))

	return ServConfig.CacheTTL, nil
}

func GetView(domain, line, area string) *View {
	return viewCache[domain][line+area]
}
