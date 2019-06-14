package main

import (
	"fmt"
	"time"

	"flag"
	"github.com/miekg/dns"
	"github.com/miekg/dns/idn"
	"strings"
)

var (
	zdnsNs = []string{"dns1.zdnscloud.info.", "ins1.zdnscloud.com.", "cns1.zdnscloud.net.", "vns1.zdnscloud.biz."}
	domain = ""
)

const (
	rootDnsIp = "198.41.0.4:53"
)

func init() {
	flag.StringVar(&domain, "d", "", "domain")
}

func compareNs(ns []string) bool {
	for _, v := range ns {
		for _, x := range zdnsNs {
			if v == x {
				return true
			}
		}
	}

	return false
}

func dig(domain, dnsServer string, digType uint16) (*dns.Msg, error) {
	if dnsServer == "" {
		dnsServer = rootDnsIp
	}
	m := new(dns.Msg)
	domain = idn.ToPunycode(domain)
	fmt.Println(fmt.Sprintf("dig @%s %s %d", dnsServer, domain, digType))
	m.SetQuestion(domain, digType)
	c := new(dns.Client)
	c.DialTimeout = 5 * time.Second
	c.ReadTimeout = 5 * time.Second
	c.WriteTimeout = 5 * time.Second
	r, _, err := c.Exchange(m, dnsServer)
	if err != nil {
		return nil, err
	}
	rcode, ok := dns.RcodeToString[r.Rcode]
	if ok && rcode == "NOERROR" {
		return r, nil
	}

	return nil, nil
}

func parseNs(r *dns.Msg) []string {
	ns := []string{}
	rrs := []dns.RR{}
	if r == nil {
		fmt.Println("message is nil")
		return ns
	}
	if len(r.Ns) != 0 {
		rrs = r.Ns
	} else if len(r.Answer) != 0 {
		rrs = r.Answer
	}
	for _, v := range rrs {
		fmt.Println(v)
		switch nstmp := v.(type) {
		case *dns.NS:
			ns = append(ns, nstmp.Ns)
		}
	}
	return ns
}

func parseExtra(ns string, r *dns.Msg) (ipv4s []string, ipv6s []string) {
	for _, v := range r.Extra {
		switch glue := v.(type) {
		case *dns.A:
			if ns == glue.Hdr.Name {
				ipv4s = append(ipv4s, glue.A.String()+":53")
			}
		case *dns.AAAA:
			if ns == glue.Hdr.Name {
				ipv4s = append(ipv4s, glue.AAAA.String()+":53")
			}
		default:
			fmt.Println("not support")
		}
	}
	return
}

func digNs(domain, dnsServer string, digType uint16) (ns []string, err error) {
	domainLabels := dns.SplitDomainName(domain)
	labelCount := dns.CountLabel(domain)

	var msg *dns.Msg

	for i := labelCount - 1; i >= 0; i-- {
		d := dns.Fqdn(strings.Join(domainLabels[i:labelCount], "."))

		msg, err = dig(d, dnsServer, dns.TypeNS)
		if err != nil {
			return
		}

		ns = parseNs(msg)
		if len(ns) == 0 {
			return
		}
		for _, v := range ns {
			ipv4s, _ := parseExtra(v, msg)
			if len(ipv4s) == 0 {
				continue
			} else {
				dnsServer = ipv4s[0]
				break
			}
		}
	}
	return
}

func main() {
	flag.Parse()
	ns, err := digNs(domain, "", 0)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println(ns)
	}
}
