// Package probe stores basic probes that are used to check services health
package probe

import (
	"fmt"
	"time"
)

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

// EvaluateDuration is a shortcut for warning duration checks.
// It returns a message containing the duration, and a OK or a WARNING status depending on the provided warning duration.
func EvaluateDuration(duration time.Duration, warning time.Duration) (status Status, message string) {
	if duration >= warning {
		status = StatusWarning
	} else {
		status = StatusOK
	}
	message = fmt.Sprintf("%d ms", duration.Nanoseconds()/1000000)
	return
}
