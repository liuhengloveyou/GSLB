package service

import (
	"fmt"

	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/geo"
)

// rr : [domain][type]RR
func ResolvDomains(clientIP string, count int, rr map[string]map[uint16][]*RR) error {
	ip, _ := geo.FindIP(clientIP)
	Logger.Info(fmt.Sprintf("ResolvDomains: %s %v", clientIP, ip))

	for domain := range rr {
		line, area := ip.GetLineArea()
		view, err := GetView(line, area)
		if err != nil {
			Logger.Error("gerview ERR: " + domain + "; " + clientIP)
			continue
		}

		// 指定view
		rrs := GetRRByView(domain, view.View)
		rrs = GroupLB(domain, view.View, ip, rrs)
		for i := 0; i < count && i < len(rrs); i++ {
			rr[rrs[i].Domain][rrs[i].Type] = append(rr[rrs[i].Domain][rrs[i].Type], rrs[i])
		}
		Logger.Info("ResolvDomains one: " + domain + "\t" + clientIP + ": " + fmt.Sprintf("%#v", view) + ": " + fmt.Sprintf("%#v", rrs))

		count = count - len(rrs)

		// 线路默认view
		if count > 0 {
			rrs = GetRRByView(domain, LineDefault(line))
			rrs = GroupLB(domain, LineDefault(line), ip, rrs)
			for i := 0; i < count && i < len(rrs); i++ {
				rr[rrs[i].Domain][rrs[i].Type] = append(rr[rrs[i].Domain][rrs[i].Type], rrs[i])
			}
			Logger.Info("ResolvDomains one: " + domain + "\t" + clientIP + ": " + LineDefault(line) + ": " + fmt.Sprintf("%#v", rrs))
		}

		count = count - len(rrs)

		// 域名默认view
		if count > 0 {
			rrs = GetRRByView(domain, ANY)
			rrs = GroupLB(domain, ANY, ip, rrs)
			for i := 0; i < count && i < len(rrs); i++ {
				rr[rrs[i].Domain][rrs[i].Type] = append(rr[rrs[i].Domain][rrs[i].Type], rrs[i])
			}
			Logger.Info("ResolvDomains one: " + domain + "\t" + clientIP + ": ANY: " + fmt.Sprintf("%#v", rrs))
		}

		/*
			} else {
				// 没有定义解析规则
				// @@@　要不要递归权威ＤＮＳ呢？
			}
		*/
	}

	Logger.Info(fmt.Sprintf("ResolvDomains OK: %#v", rr))

	return nil
}
