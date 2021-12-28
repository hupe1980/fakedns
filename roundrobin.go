package fakedns

import (
	"net"
	"sync/atomic"
)

type RoundRobin interface {
	Next() net.IP
	HasEntries() bool
}

type roundrobin struct {
	ips  []net.IP
	next uint32
}

func NewRoundRobin(ips ...string) RoundRobin {
	parsedIPs := []net.IP{}
	for _, ip := range ips {
		parsedIPs = append(parsedIPs, net.ParseIP(ip))
	}

	return &roundrobin{
		ips: parsedIPs,
	}
}

func (r *roundrobin) HasEntries() bool {
	return len(r.ips) > 0
}

// Next returns next ip
func (r *roundrobin) Next() net.IP {
	n := atomic.AddUint32(&r.next, 1)
	return r.ips[(int(n)-1)%len(r.ips)]
}
