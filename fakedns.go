package fakedns

import (
	"crypto/tls"
	"net"
	"regexp"
	"strings"

	"github.com/miekg/dns"
)

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
	server  *dns.Server
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

	server := &dns.Server{
		Handler: &handler{
			rebind:      options.Rebind,
			ttl:         options.TTL,
			fallbackDNS: options.FallbackDNSResolver,
			re:          re,
			ipV4Pool:    NewRoundRobin(options.IPsV4...),
			ipV6Pool:    NewRoundRobin(options.IPsV6...),
			text:        options.Text,
		},
	}

	fakeDNS := &FakeDNS{
		options: options,
		server:  server,
	}

	return fakeDNS, nil
}

func (t *FakeDNS) ListenAndServe(addr, network string) error {
	if addr != "" {
		if _, _, err := net.SplitHostPort(addr); err != nil {
			addr = net.JoinHostPort(addr, "53")
		}
	}

	t.server.Addr = addr
	t.server.Net = network

	return t.server.ListenAndServe()
}

func (t *FakeDNS) ListenAndServeTLS(addr, certFile, keyFile string) error {
	if addr != "" {
		if _, _, err := net.SplitHostPort(addr); err != nil {
			addr = net.JoinHostPort(addr, "53")
		}
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	t.server.TLSConfig = config
	t.server.Addr = addr
	t.server.Net = "tcp-tls"

	return t.server.ListenAndServe()
}

func (t *FakeDNS) Shutdown() error {
	return t.server.Shutdown()
}
