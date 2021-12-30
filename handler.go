package fakedns

import (
	"net"
	"regexp"
	"strings"

	"github.com/miekg/dns"
)

type handler struct {
	ttl         uint32
	fallbackDNS string
	rebind      *Rebind
	re          *regexp.Regexp
	ipV4Pool    RoundRobin
	ipV6Pool    RoundRobin
	text        []string
	logger      Logger
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := &dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, q := range r.Question {
		h.logger.Printf(INFO, "[*] Receiving question: %s", q.String())

		domain := r.Question[0].Name
		domainlookup := strings.TrimSuffix(domain, ".")

		if h.rebind != nil {
			h.rebind.Inc(domainlookup)
		}

		if h.re.MatchString(domainlookup) {
			switch q.Qtype {
			case dns.TypeA:
				if h.ipV4Pool.HasEntries() {
					msg.Answer = append(msg.Answer, &dns.A{
						Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: h.ttl},
						A:   h.ipV4(domainlookup),
					})
				}
			case dns.TypeAAAA:
				if h.ipV6Pool.HasEntries() {
					msg.Answer = append(msg.Answer, &dns.AAAA{
						Hdr:  dns.RR_Header{Name: domain, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: h.ttl},
						AAAA: h.ipV6(domainlookup),
					})
				}
			case dns.TypeTXT:
				if len(h.text) > 0 {
					msg.Answer = append(msg.Answer, &dns.TXT{
						Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: h.ttl},
						Txt: h.text,
					})
				}
			}
		} else if h.fallbackDNS != "" {
			if exMsg, err := dns.Exchange(r, h.fallbackDNS); err == nil {
				msg.Answer = append(msg.Answer, exMsg.Answer...)
			}
		}
	}

	err := w.WriteMsg(msg)
	if err != nil {
		h.logger.Printf(ERROR, "[!] Cannot write msg: %s", err)
	}
}

func (h *handler) ipV4(domain string) net.IP {
	if h.rebind != nil && h.rebind.IsV4Activ(domain) {
		return h.rebind.IPV4()
	}

	return h.ipV4Pool.Next()
}

func (h *handler) ipV6(domain string) net.IP {
	if h.rebind != nil && h.rebind.IsV6Activ(domain) {
		return h.rebind.IPV6()
	}

	return h.ipV6Pool.Next()
}
