package probe

import (
	"github.com/miekg/dns"
	"github.com/mitchellh/mapstructure"
)

// DNS Probe, used to check whether a DNS server is answering.
// `domain` will be resolved through a lookup for an A record.
// `expected` should be the first returned IPv4 address or empty to accept any IP address.
type DNS struct {
	Config
	domain, expected string
}

// Init configures the probe.
func (d *DNS) Init(c Config) error {
	err := mapstructure.Decode(c.Options, d)
	d.Config = c
	return err
}

// Probe checks a DNS server.
// If the operation succeeds, the message will be the duration of the dial in ms.
// Otherwise, an error message is returned.
func (d *DNS) Probe() (status Status, message string) {
	m := new(dns.Msg)
	m.SetQuestion(d.domain, dns.TypeA)

	c := new(dns.Client)
	r, rtt, err := c.Exchange(m, d.Target+":53")
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

	return EvaluateDuration(rtt, d.Warning)
}
