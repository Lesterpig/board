package probe

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// SMTP Probe, used to check smtp servers status
// BEWARE! Only full TLS servers are working with this probe.
type SMTP struct {
	ProberConfig
}

// Init configures the probe.
func (s *SMTP) Init(c ProberConfig) error {
	s.ProberConfig = c
	return nil
}

// Probe checks a mailbox status.
// If the operation succeeds, the message will be the duration of the SMTP handshake in ms.
// Otherwise, an error message is returned.
func (s *SMTP) Probe() (status Status, message string) {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", s.Target, s.Fatal)
	if err != nil {
		return StatusError, defaultConnectErrorMsg
	}

	defer func() { _ = conn.Close() }()

	host, _, _ := net.SplitHostPort(s.Target)
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

	return EvaluateDuration(time.Since(start), s.Warning)
}
