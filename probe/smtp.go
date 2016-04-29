package probe

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// SMTP Probe, used to check smtp servers status
type SMTP struct {
	addrport       string
	warning, fatal time.Duration
}

// NewSMTP returns a ready-to-go probe.
// A warning will be triggered if the response takes more than `warning` to come.
// BEWARE! Only full TLS servers are working with this probe.
func NewSMTP(addrport string, warning, fatal time.Duration) *SMTP {
	return &SMTP{
		addrport: addrport,
		warning:  warning,
		fatal:    fatal,
	}
}

// Probe checks a website status.
// If the operation succeeds, the message will be the duration of the SMTP handshake in ms.
// Otherwise, an error message is returned.
func (s *SMTP) Probe() (status Status, message string) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", s.addrport, s.fatal)
	if err != nil {
		return StatusError, "Unable to connect"
	}

	defer func() { _ = conn.Close() }()
	host, _, _ := net.SplitHostPort(s.addrport)
	secure := tls.Client(conn, &tls.Config{
		ServerName: host,
	})

	data := make([]byte, 4)
	_, err = secure.Read(data)
	if err != nil {
		return StatusError, "TLS Error"
	}
	if fmt.Sprintf("%s", data) != "220 " {
		return StatusError, "Unexpected reply"
	}

	return EvaluateDuration(time.Since(start), s.warning)
}
