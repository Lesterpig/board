package probe

import (
	"time"

	ping "github.com/sparrc/go-ping"
)

const packetCount = 3

// Ping probe, check network avaibility
type Ping struct {
	addr    string
	warning time.Duration
	fail    time.Duration
}

// NewPing returns a ping probe
func NewPing(addr string, warning, fail time.Duration) *Ping {
	return &Ping{
		addr:    addr,
		warning: warning,
		fail:    fail,
	}
}

// Probe send 3 ICMP packets and wait for replies
// Warning is issued when a reply is missing or if rtt exeed `warning`
// Failure occure when target is unreachable
func (p *Ping) Probe() (status Status, message string) {
	pinger := &ping.Pinger{
		Timeout: p.fail,
		Count:   packetCount,
	}
	pinger.SetAddr(p.addr)

	pinger.Run()
	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		return StatusError, "Unreacheable"
	}

	if stats.PacketsRecv < packetCount {
		return StatusWarning, "Packet loss"
	}
	return EvaluateDuration(stats.AvgRtt, p.warning)
}
