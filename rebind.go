package fakedns

import (
	"net"
)

type Rebind struct {
	rebindV4  string
	rebindV6  string
	threshold uint32
	counter   map[string]uint32
}

func NewRebind(rebindV4, rebindV6 string, threshold uint32) *Rebind {
	return &Rebind{
		rebindV4:  rebindV4,
		rebindV6:  rebindV6,
		threshold: threshold,
		counter:   make(map[string]uint32),
	}
}

func (r *Rebind) Inc(domain string) {
	r.counter[domain]++
}

func (r *Rebind) IsV4Activ(domain string) bool {
	if r.rebindV4 == "" {
		return false
	}

	if c, ok := r.counter[domain]; ok {
		return c > r.threshold
	}

	return false
}

func (r *Rebind) IsV6Activ(domain string) bool {
	if r.rebindV6 == "" {
		return false
	}

	if c, ok := r.counter[domain]; ok {
		return c > r.threshold
	}

	return false
}

func (r *Rebind) IPV4() net.IP {
	return net.ParseIP(r.rebindV4)
}

func (r *Rebind) IPV6() net.IP {
	return net.ParseIP(r.rebindV6)
}
