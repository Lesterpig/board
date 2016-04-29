// Package probe stores basic probes that are used to check services health
package probe

// Status represents the current status of a monitored service.
type Status string

// These constants represent the different available statuses of a service.
const (
	StatusUnknown Status = ""
	StatusWarning        = "WARNING"
	StatusError          = "ERROR"
	StatusOK             = "OK"
)

// Prober is the base interface that each probe must implement.
type Prober interface {
	// Probe is expected to check one service's health.
	// An additionnal message can be returned for more feedbacks.
	Probe() (status Status, message string)
}
