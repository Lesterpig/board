package probe

import (
	"github.com/miekg/dns"
	"time"
)

// DNS Probe, used to check whether a DNS server is answering.
type DNS struct {
	addr, domain, expected string
	warning, fatal         time.Duration
}

// NewDNS returns a ready-to-go probe.
// `domain` will be resolved through a lookup for an A record.
// `expected` should be the first returned IPv4 address or empty to accept any IP address.
// A warning will be triggered if the response takes more than `warning` to come.
func NewDNS(addr, domain, expected string, warning, fatal time.Duration) *DNS {
	return &DNS{
		addr:     addr,
		domain:   domain,
		expected: expected,
		warning:  warning,
		fatal:    fatal,
	}
}

// Probe checks a DNS server.
// If the operation succeeds, the message will be the duration of the dial in ms.
// Otherwise, an error message is returned.
func (d *DNS) Probe() (status Status, message string) {
	m := new(dns.Msg)
	m.SetQuestion(d.domain, dns.TypeA)

	c := new(dns.Client)
	r, rtt, err := c.Exchange(m, d.addr+":53")
	if err != nil {
		return StatusError, err.Error()
	}

	if r.Rcode != dns.RcodeSuccess {
		return StatusError, "Failed to resolve domain."
	}

	if answer, ok := r.Answer[0].(*dns.A); ok {
		if d.expected != "" && answer.A.String() != d.expected {
			return StatusError, "Unexpected DNS answer"
		}
	} else {
		return StatusError, "Failed to resolve domain."
	}

	return EvaluateDuration(rtt, d.warning)
}
