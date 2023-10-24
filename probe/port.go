package probe

import (
	"net"
	"net/url"
	"time"
)

// Port Probe, used to check if a port is open or not.
type Port struct {
	ProberConfig
	network, addrport string
}

// Init configures the probe.
func (p *Port) Init(c ProberConfig) error {
	p.ProberConfig = c

	u, err := url.Parse(p.Target)
	if err != nil {
		return err
	}

	p.network = u.Scheme
	p.addrport = u.Host

	return nil
}

// Probe checks a port status
// If the operation succeeds, the message will be the duration of the dial in ms.
// Otherwise, an error message is returned.
func (p *Port) Probe() (status Status, message string) {
	start := time.Now()

	conn, err := net.DialTimeout(p.network, p.addrport, p.Fatal)
	if err != nil {
		return StatusError, defaultConnectErrorMsg
	}

	_ = conn.Close()

	return EvaluateDuration(time.Since(start), p.Warning)
}
