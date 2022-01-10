package fakedns

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type Rebind struct {
	rebindV4  net.IP
	rebindV6  net.IP
	threshold int
	counter   map[string]int
	mu        sync.RWMutex
}

func NewRebind(rebindV4, rebindV6 string, threshold int) (*Rebind, error) {
	var ipv4, ipv6 net.IP

	if rebindV4 != "" {
		ipv4 = net.ParseIP(rebindV4)

		if strings.Contains(rebindV4, ":") || ipv4 == nil {
			return nil, fmt.Errorf("invalid IPV4 address: %s", rebindV4)
		}
	}

	if rebindV6 != "" {
		ipv6 = net.ParseIP(rebindV6)

		if strings.Contains(rebindV6, ".") || ipv6 == nil {
			return nil, fmt.Errorf("invalid IPV6 address: %s", rebindV6)
		}
	}

	return &Rebind{
		rebindV4:  ipv4,
		rebindV6:  ipv6,
		threshold: threshold,
		counter:   make(map[string]int),
	}, nil
}

func (r *Rebind) Inc(domain string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter[domain]++
}

func (r *Rebind) IsV4Activ(domain string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.rebindV4 == nil {
		return false
	}

	if c, ok := r.counter[domain]; ok {
		return c > r.threshold
	}

	return false
}

func (r *Rebind) IsV6Activ(domain string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.rebindV6 == nil {
		return false
	}

	if c, ok := r.counter[domain]; ok {
		return c > r.threshold
	}

	return false
}

func (r *Rebind) IPV4() net.IP {
	return r.rebindV4
}

func (r *Rebind) IPV6() net.IP {
	return r.rebindV6
}
