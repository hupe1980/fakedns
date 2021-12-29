package fakedns

import (
	"net"
	"regexp"
	"strings"

	"github.com/miekg/dns"
)

type fakeDNSHandler struct {
	ttl         uint32
	fallbackDNS string
	rebind      *Rebind
	re          *regexp.Regexp
	ipV4Pool    RoundRobin
	ipV6Pool    RoundRobin
	text        []string
}

func (t *fakeDNSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	domain := r.Question[0].Name
	domainlookup := strings.TrimSuffix(domain, ".")

	if t.rebind != nil {
		t.rebind.Inc(domainlookup)
	}

	msg := &dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	if t.re.MatchString(domainlookup) {
		for _, q := range r.Question {
			switch q.Qtype {
			case dns.TypeA:
				if t.ipV4Pool.HasEntries() {
					msg.Answer = append(msg.Answer, &dns.A{
						Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: t.ttl},
						A:   t.ipV4(domainlookup),
					})
				}
			case dns.TypeAAAA:
				if t.ipV6Pool.HasEntries() {
					msg.Answer = append(msg.Answer, &dns.AAAA{
						Hdr:  dns.RR_Header{Name: domain, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: t.ttl},
						AAAA: t.ipV6(domainlookup),
					})
				}
			case dns.TypeTXT:
				if len(t.text) > 0 {
					msg.Answer = append(msg.Answer, &dns.TXT{
						Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: t.ttl},
						Txt: t.text,
					})
				}
			}

		}
	} else if t.fallbackDNS != "" {
		if exMsg, err := dns.Exchange(r, t.fallbackDNS); err == nil {
			msg.Answer = append(msg.Answer, exMsg.Answer...)
		}
	}

	_ = w.WriteMsg(msg)
}

func (t *fakeDNSHandler) ipV4(domain string) net.IP {
	if t.rebind != nil && t.rebind.IsV4Activ(domain) {
		return t.rebind.IPV4()
	}

	return t.ipV4Pool.Next()
}

func (t *fakeDNSHandler) ipV6(domain string) net.IP {
	if t.rebind != nil && t.rebind.IsV6Activ(domain) {
		return t.rebind.IPV6()
	}

	return t.ipV6Pool.Next()
}

type Options struct {
	FallbackDNSResolver string
	TTL                 uint32
	Domains             []string
	IPsV4               []string
	IPsV6               []string
	Rebind              *Rebind
	Text                []string
}

type FakeDNS struct {
	options *Options
	handler dns.Handler
}

func (t *FakeDNS) ListenAndServe(addr, network string) error {
	if addr != "" {
		if _, _, err := net.SplitHostPort(addr); err != nil {
			addr = net.JoinHostPort(addr, "53")
		}
	}

	srv := &dns.Server{
		Addr:    addr,
		Net:     network,
		Handler: t.handler,
	}

	return srv.ListenAndServe()
}

func New(options *Options) (*FakeDNS, error) {
	re, err := regexp.Compile(strings.Join(options.Domains, "|"))
	if err != nil {
		return nil, err
	}

	if options.FallbackDNSResolver != "" {
		if _, _, err := net.SplitHostPort(options.FallbackDNSResolver); err != nil {
			options.FallbackDNSResolver = net.JoinHostPort(options.FallbackDNSResolver, "53")
		}
	}

	fakeDNS := &FakeDNS{
		options: options,
		handler: &fakeDNSHandler{
			rebind:      options.Rebind,
			ttl:         options.TTL,
			fallbackDNS: options.FallbackDNSResolver,
			re:          re,
			ipV4Pool:    NewRoundRobin(options.IPsV4...),
			ipV6Pool:    NewRoundRobin(options.IPsV6...),
			text:        options.Text,
		},
	}

	return fakeDNS, nil
}
