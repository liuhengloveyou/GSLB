package api

import (
	"fmt"
	"testing"

	. "github.com/miekg/dns"
)

func HelloServerEchoAddrPort(w ResponseWriter, req *Msg) {
	m := new(Msg)
	m.SetReply(req)

	remoteAddr := w.RemoteAddr().String()
	m.Extra = make([]RR, 1)
	m.Extra[0] = &TXT{Hdr: RR_Header{Name: m.Question[0].Name, Rrtype: TypeTXT, Class: ClassINET, Ttl: 0}, Txt: []string{remoteAddr}}
	w.WriteMsg(m)
}

func AnotherHelloServer(w ResponseWriter, req *Msg) {
	m := new(Msg)
	m.SetReply(req)

	m.Extra = make([]RR, 1)
	m.Extra[0] = &TXT{Hdr: RR_Header{Name: m.Question[0].Name, Rrtype: TypeTXT, Class: ClassINET, Ttl: 0}, Txt: []string{"Hello example"}}
	w.WriteMsg(m)
}

func TestServing(t *testing.T) {
	InitDnsApi()
}

func TestClient(t *testing.T) {
	addrstr := "[::]:1320"
	c := new(Client)
	m := new(Msg)
	m.SetQuestion("miek.nl.", TypeTXT)
	r, _, err := c.Exchange(m, addrstr)
	if err != nil || len(r.Extra) == 0 {
		t.Fatal("failed to exchange miek.nl", err)
	}
	txt := r.Extra[0].(*TXT).Txt[0]
	if txt != "Hello world" {
		t.Error("unexpected result for miek.nl", txt, "!= Hello world")
	}

	m.SetQuestion("example.com.", TypeTXT)
	r, _, err = c.Exchange(m, addrstr)
	if err != nil {
		t.Fatal("failed to exchange example.com", err)
	}
	txt = r.Extra[0].(*TXT).Txt[0]
	if txt != "Hello example" {
		t.Error("unexpected result for example.com", txt, "!= Hello example")
	}
	t.Log(txt)
	fmt.Println(">>>>>>>>>>>>>>", txt)

	// Test Mixes cased as noticed by Ask.
	m.SetQuestion("eXaMplE.cOm.", TypeTXT)
	r, _, err = c.Exchange(m, addrstr)
	if err != nil {
		t.Error("failed to exchange eXaMplE.cOm", err)
	}
	txt = r.Extra[0].(*TXT).Txt[0]
	if txt != "Hello example" {
		t.Error("unexpected result for example.com", txt, "!= Hello example")
	}
	t.Log(txt)
}

/*nc TestDNS(t *testing.T) {
	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	c := new(dns.Client)
	m := new(dns.Msg)
	zone := "10zan.net"
	m.SetQuestion(dns.Fqdn(zone), dns.TypeANY)
	m.SetEdns0(4096, true)

	r, _, err := c.Exchange(m, config.Servers[0]+":"+config.Port)
	if err != nil {
		return
	}
	if r.Rcode != dns.RcodeSuccess {
		return
	}
	for _, k := range r.Answer {
		fmt.Println(">>>", k)
		if key, ok := k.(*dns.DNSKEY); ok {
			for _, alg := range []uint8{dns.SHA1, dns.SHA256, dns.SHA384} {
				fmt.Printf("%s; %d\n", key.ToDS(alg).String(), key.Flags)
			}
		} else if a, ok := k.(*dns.A); ok {
			fmt.Printf("%s\n", a.String())
		}
	}
}
*/
