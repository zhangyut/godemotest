package main

import (
	"errors"
	"fmt"

	"flag"
	dns "github.com/miekg/dns"
	"strings"
)

var (
	domain string
)

func init() {
	flag.StringVar(&domain, "d", "", "domain name")
}

func findNs(domain string) (error, bool) {
	zdnsNs := []string{"dns1.zdnscloud.info.", "ins1.zdnscloud.com.", "cns1.zdnscloud.net.", "vns1.zdnscloud.biz."}
	domainName := domain
	domainLabels := dns.SplitDomainName(domainName)
	labelCount := dns.CountLabel(domainName)

	dnsServer := "198.41.0.4:53"

	for i := labelCount - 1; i > 0; i-- {
		domain := dns.Fqdn(strings.Join(domainLabels[i:labelCount], "."))
		m := new(dns.Msg)
		m.SetQuestion(domain, dns.TypeNS)
		fmt.Println("+dig " + dnsServer + " " + domain)
		c := new(dns.Client)
		r, _, err := c.Exchange(m, dnsServer)
		if err != nil {
			return err, false
		}

		for _, ns := range r.Ns {
			switch nstmp := ns.(type) {
			case *dns.NS:
				for _, glue := range r.Extra {
					switch atmp := glue.(type) {
					case *dns.A:
						if nstmp.Ns == atmp.Hdr.Name {
							dnsServer = atmp.A.String() + ":53"
							if i == 1 {
								goto DIGNS
							} else {
								continue
							}
						}
					}
				}
			}
		}
	}

DIGNS:
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domainName), dns.TypeNS)
	c := new(dns.Client)
	fmt.Println("-dig " + dnsServer + " " + domain)
	r, _, err := c.Exchange(m, dnsServer)
	if err != nil {
		return err, false
	}
	if r.MsgHdr.Rcode != dns.RcodeSuccess {
		return errors.New(fmt.Sprintf("error rcode:%d", r.MsgHdr.Rcode)), false
	}
	hasNs := false
	for _, ns := range r.Ns {
		fmt.Println("ns: " + dns.Fqdn(strings.ToLower(ns.(*dns.NS).Ns)))
		for _, zns := range zdnsNs {
			if dns.Fqdn(strings.ToLower(ns.(*dns.NS).Ns)) == zns {
				hasNs = true
			}
		}

	}
	if !hasNs {
		return errors.New("not match ns"), false
	}

	return nil, true
}

func main() {
	flag.Parse()
	if domain == "" {
		fmt.Println("domain name is nil.")
		return
	}
	err, ok := findNs(domain)
	if err == nil && ok {
		fmt.Println("ok.")
	} else {
		fmt.Println("false.")
	}
}
