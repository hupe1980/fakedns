package fakedns

import (
	"crypto/tls"
	"log"
	"net"
	"regexp"

	"github.com/miekg/dns"
)

type Options struct {
	FallbackDNSResolver string
	TTL                 uint32
	IPsV4               []string
	IPsV6               []string
	Rebind              *Rebind
	Text                []string
	MX                  string
	Logger              Logger
}

type FakeDNS struct {
	options *Options
	server  *dns.Server
	logger  Logger
}

func New(domain *regexp.Regexp, options *Options) *FakeDNS {
	if options.FallbackDNSResolver != "" {
		if _, _, err := net.SplitHostPort(options.FallbackDNSResolver); err != nil {
			options.FallbackDNSResolver = net.JoinHostPort(options.FallbackDNSResolver, "53")
		}
	}

	if options.Logger == nil {
		options.Logger = NewDefaultLogger(INFO, log.Default())
	}

	server := &dns.Server{
		Handler: &handler{
			rebind:      options.Rebind,
			ttl:         options.TTL,
			fallbackDNS: options.FallbackDNSResolver,
			re:          domain,
			ipV4Pool:    NewRoundRobin(options.IPsV4...),
			ipV6Pool:    NewRoundRobin(options.IPsV6...),
			text:        options.Text,
			mx:          options.MX,
			logger:      options.Logger,
		},
	}

	fakeDNS := &FakeDNS{
		options: options,
		server:  server,
		logger:  options.Logger,
	}

	return fakeDNS
}

func (t *FakeDNS) ListenAndServe(addr, network string) error {
	if addr != "" {
		if _, _, err := net.SplitHostPort(addr); err != nil {
			addr = net.JoinHostPort(addr, "53")
		}
	}

	t.server.Addr = addr
	t.server.Net = network

	t.logger.Printf(INFO, "[*] Starting server on %s", addr)

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
		MinVersion:   tls.VersionTLS12,
	}

	t.server.TLSConfig = config
	t.server.Addr = addr
	t.server.Net = "tcp-tls"

	t.logger.Printf(INFO, "Starting server on %s", addr)

	return t.server.ListenAndServe()
}

func (t *FakeDNS) Shutdown() error {
	return t.server.Shutdown()
}
