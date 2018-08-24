package service

import (
	"fmt"
	"time"

	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/dao"
	//	"github.com/liuhengloveyou/GSLB/geo"
)

// 域名记录缓存
// [domain][group]RR
var rrCache = make(map[string]map[string]*RR)

// 解析规则缓存
var cache = make(map[string]map[string]interface{})

// [domain][type]RR
func ResolvDomains(clientIP string, rr map[string]map[uint16]*RR) error {
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
		//	line, area := geo.FindIP(clientIP)

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

// SELECT rule.domain, zone.line, zone.area,rule.group FROM rule join zone on rule.zone = zone.zone;

func LoadRRCache() {
	now := "0"

	for {
		sleep := ServConfig.CacheTTL

		rr, err := dao.LoadRRFromMysql(now)
		if err != nil {
			Logger.Error("LoadRRFromMysql ERR: " + err.Error())
			time.Sleep(time.Second * time.Duration(ServConfig.CacheTTL))
		}

		for _, r := range rr {
			if int64(r.Ttl) < sleep {
				sleep = int64(r.Ttl)
			}
			fmt.Println(">>>", r.Domain, r.Group)
		}

		now = time.Now().Format("2006-01-02 15:04:05")
		time.Sleep(time.Second * time.Duration(sleep))
	}
}
