package service

import (
	"fmt"

	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/dao"
)

// [domain][type]RR
func ResolvDomains(client string, rr map[string]map[uint16]*RR) error {

	// 所有域名
	d := make([]string, len(rr))
	i := 0
	for k := range rr {
		d[i] = k
		i++
	}

	// 查就近解析规则
	rules, err := dao.SelectRulesFromMysql(d)
	if err != nil {
		Logger.Error("dao.SelectRulesFromMysql ERR: " + err.Error())
		return err
	}

	// 有定义就近解析规则就按规则解析
	if len(rules) > 0 {

	} else {
		// 没有定义解析规则
	}

	// 查询域名资源记录
	r, e := dao.SelectRRsFromMysql(d)
	if e != nil {
		return e
	}

	for i := 0; i < len(r); i++ {
		domain := r[i].Domain
		rtype := r[i].Type
		Logger.Debug("resolv one: " + domain + string(rtype))

		d := rr[domain]
		if _, ok := d[rtype]; ok {
			rr[domain][rtype] = r[i]
		}
	}

	Logger.Info(fmt.Sprintf("Resolved: %#v", rr))

	return nil
}
