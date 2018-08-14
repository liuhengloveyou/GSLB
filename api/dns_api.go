// dig -b 127.0.0.1 -p 1053
package api

import (
	"net"
	"time"

	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/service"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

func rootDNServer(w dns.ResponseWriter, req *dns.Msg) {
	qq := make(map[string]map[uint16]*RR)
	for _, q := range req.Question {
		log.Infoln("DNS question:", q.Name, q.Qtype)
		if qt, ok := qq[q.Name]; ok {
			qt[q.Qtype] = nil
		} else {
			qq[q.Name] = map[uint16]*RR{q.Qtype: nil}
		}
	}

	if err := service.ResolvDomains(qq); err != nil {
		log.Errorln("DNS resolv ERR: ", err)
		return
	}

	m := new(dns.Msg)
	m.SetReply(req)

	for domain, v := range qq {
		for rtype, rr := range v {
			if rr == nil {
				continue
			}

			switch rtype {
			case dns.TypeA:
				m.Answer = append(m.Answer, &dns.A{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    rr.Ttl,
					},
					A: net.ParseIP(rr.Data),
				})
			case dns.TypeCNAME:
				m.Answer = append(m.Answer, &dns.CNAME{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeCNAME,
						Class:  dns.ClassINET,
						Ttl:    rr.Ttl,
					},
					Target: rr.Data,
				})
			}
		}
	}

	err := w.WriteMsg(m)
	if err != nil {
		log.Errorln("DNSRootServer ERR:", err)
		return
	}

	log.Infoln("DNSRootServer OK:", m)
	return
}

func InitDnsApi(addr string) error {

	dns.HandleFunc(".", rootDNServer)

	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}

	server := &dns.Server{PacketConn: pc, ReadTimeout: time.Minute, WriteTimeout: time.Minute}
	if err = server.ActivateAndServe(); err != nil {
		return err
	}

	return nil
}
