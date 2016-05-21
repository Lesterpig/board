package probe

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"crypto/tls"
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
// Set `verifyCertificate` to `false` to skip the certificate verification.
func NewHTTP(addrport string, warning, fatal time.Duration, regex string, verifyCertificate bool) *HTTP {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifyCertificate},
	}
	return &HTTP{
		client: &http.Client{
			Timeout:   fatal,
			Transport: tr,
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

	return EvaluateDuration(duration, h.warning)
}
