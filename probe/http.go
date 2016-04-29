package probe

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// HTTP Probe, used to check HTTP(S) websites status.
type HTTP struct {
	client   *http.Client
	addrport string
	warning  time.Duration
	regex    *regexp.Regexp
}

// NewHTTP returns a ready-to-go probe.
// A warning will be triggered if the response takes more than `warning` to come.
// The `regex` is used to check the content of the website, and can be empty.
func NewHTTP(addrport string, warning time.Duration, fatal time.Duration, regex string) *HTTP {
	return &HTTP{
		client: &http.Client{
			Timeout: fatal,
		},
		addrport: addrport,
		warning:  warning,
		regex:    regexp.MustCompile(regex),
	}
}

// Probe checks a website status.
// If the operation succeeds, the message will be the duration of the HTTP request in ms.
// Otherwise, an error message is returned.
func (h *HTTP) Probe() (status Status, message string) {
	start := time.Now()
	res, err := h.client.Get(h.addrport)
	duration := time.Since(start)

	if err != nil {
		return StatusError, "Unable to connect"
	}

	if res.StatusCode != 200 {
		return StatusError, strconv.Itoa(res.StatusCode)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if !h.regex.Match(body) {
		return StatusError, "Unexpected result"
	}

	if duration >= h.warning {
		status = StatusWarning
	} else {
		status = StatusOK
	}

	message = fmt.Sprintf("%d ms", duration.Nanoseconds()/1000000)
	return
}
