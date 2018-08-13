package api

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

func RootServer(w dns.ResponseWriter, req *dns.Msg) {

	fmt.Printf("%#v\n", req)

	m := new(dns.Msg)
	m.SetReply(req)

	m.Answer = make([]dns.RR, 1)
	m.Answer[0] = &dns.A{Hdr: dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}, A: net.ParseIP("1.1.1.1")}
	m.Extra = make([]dns.RR, 1)
	m.Extra[0] = &dns.TXT{Hdr: dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0}, Txt: []string{"Hello world"}}
	w.WriteMsg(m)
}

func InitDnsApi() {

	dns.HandleFunc(".", RootServer)

	pc, err := net.ListenPacket("udp", ":53")
	if err != nil {
		panic(err)
	}

	fmt.Println(pc.LocalAddr().String())

	server := &dns.Server{PacketConn: pc, ReadTimeout: time.Minute, WriteTimeout: time.Minute}

	if err = server.ActivateAndServe(); err != nil {
		panic(err)
	}
}
