package service

import (
	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/dao"
	log "github.com/sirupsen/logrus"
)

// [domain][type]RR
func ResolvDomains(rr map[string]map[uint16]*RR) error {
	d := make([]string, len(rr))
	i := 0
	for k := range rr {
		d[i] = k
		i++
	}

	r, e := dao.SelectRRsFromMysql(d)
	if e != nil {
		return e
	}

	for i := 0; i < len(r); i++ {
		domain := r[i].Domain
		rtype := r[i].Type
		log.Debugln("resolv one: ", domain, rtype)

		d := rr[domain]
		if _, ok := d[rtype]; ok {
			rr[domain][rtype] = r[i]
		}
	}

	log.Infof("Resolved: %#v", rr)

	return nil
}
