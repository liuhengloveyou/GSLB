// dig www.g.com @127.0.0.1 -p 1053

package api

import (
	"fmt"
	"net"
	"time"

	. "github.com/liuhengloveyou/GSLB/common"
	"github.com/liuhengloveyou/GSLB/service"

	"github.com/miekg/dns"
	"go.uber.org/zap"
)

func rootDNServer(w dns.ResponseWriter, req *dns.Msg) {
	qq := make(map[string]map[uint16][]*RR)
	for _, q := range req.Question {
		Logger.Info("DNS question:", zap.String("name", q.Name), zap.Uint16("qtype", q.Qtype))
		if qt, ok := qq[q.Name]; ok {
			qt[q.Qtype] = nil
		} else {
			qq[q.Name] = map[uint16][]*RR{q.Qtype: nil}
		}
	}

	if err := service.ResolvDomains(w.RemoteAddr().(*net.UDPAddr).String(), 1, qq); err != nil {
		Logger.Error("DNS resolv ERR: " + err.Error())
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
						Ttl:    rr[0].TTL,
					},
					A: net.ParseIP(rr[0].Record),
				})
			case dns.TypeCNAME:
				m.Answer = append(m.Answer, &dns.CNAME{
					Hdr: dns.RR_Header{
						Name:   domain,
						Rrtype: dns.TypeCNAME,
						Class:  dns.ClassINET,
						Ttl:    rr[0].TTL,
					},
					Target: rr[0].Record,
				})
			}
		}
	}

	err := w.WriteMsg(m)
	if err != nil {
		Logger.Error("DNSRootServer ERR:" + err.Error())
		return
	}

	Logger.Info(fmt.Sprintf("DNSRootServer OK: %#v", m))
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
