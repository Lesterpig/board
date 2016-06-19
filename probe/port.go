package probe

import (
	"net"
	"time"
)

// Port Probe, used to check if a port is open or not.
type Port struct {
	addrport, network string
	warning, fatal    time.Duration
}

// NewPort returns a ready-to-go probe.
// The `network` variable should be `tcp` or `udp` or their v4/v6 variants.
// A warning will be triggered if the response takes more than `warning` to come.
func NewPort(network, addrport string, warning, fatal time.Duration) *Port {
	return &Port{
		network:  network,
		addrport: addrport,
		warning:  warning,
		fatal:    fatal,
	}
}

// Probe checks a port status
// If the operation succeeds, the message will be the duration of the dial in ms.
// Otherwise, an error message is returned.
func (p *Port) Probe() (status Status, message string) {
	start := time.Now()
	conn, err := net.DialTimeout(p.network, p.addrport, p.fatal)
	if err != nil {
		return StatusError, defaultConnectErrorMsg
	}

	_ = conn.Close()
	return EvaluateDuration(time.Since(start), p.warning)
}
