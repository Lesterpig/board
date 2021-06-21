package probe

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"regexp"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

// HTTP Probe, used to check HTTP(S) websites status.
type HTTP struct {
	Config
	client *http.Client
	regex  *regexp.Regexp
}

// HTTPOptions is a structure containing optional parameters.
// The `Regex` is used to check the content of the website, and can be empty.
// Set `VerifyCertificate` to `false` to skip the TLS certificate verification.
type HTTPOptions struct {
	Regex             string
	VerifyCertificate bool
}

// Init configures the probe.
func (h *HTTP) Init(c Config) error {
	h.Config = c

	var opts HTTPOptions

	err := mapstructure.Decode(c.Options, &opts)
	if err != nil {
		return err
	}

	/* #nosec G402 */
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// Disable keep alive to get more consistent measurements
	tr.DisableKeepAlives = true
	tr.ResponseHeaderTimeout = c.Fatal

	h.client = &http.Client{
		Timeout:   c.Fatal,
		Transport: tr,
	}
	h.regex, err = regexp.Compile(opts.Regex)

	return err
}

// Probe checks a website status.
// If the operation succeeds, the message will be the duration of the HTTP request in ms.
// Otherwise, an error message is returned.
func (h *HTTP) Probe() (status Status, message string) {

	req, _ := http.NewRequest("GET", h.Target, nil)

	var server, start time.Time
	var initDuration, serverDuration, totalDuration time.Duration

	trace := &httptrace.ClientTrace{

		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			server = time.Now()
			initDuration = time.Since(start)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	// Starting timer
	start = time.Now()
	// Set server time before request. This may be reset in TLSHandshakeDone if a TLSHandshake is preformed
	server = start

	res, err := h.client.Transport.RoundTrip(req)
	totalDuration = time.Since(start)
	serverDuration = time.Since(server)

	if res != nil {
		defer res.Body.Close() // MUST CLOSED THIS
	}

	if err != nil {
		return StatusError, defaultConnectErrorMsg
	}

	if !(res.StatusCode >= 200 && res.StatusCode <= 399) {
		return StatusError, strconv.Itoa(res.StatusCode)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if h.regex != nil && !h.regex.Match(body) {
		return StatusError, "Unexpected result"
	}

	return AdvancedEvaluateDuration(initDuration, serverDuration, totalDuration, h.Warning)
}

// AdvancedEvaluateDuration is a shortcut for warning duration checks.
// It returns a message containing the duration, and a OK or a WARNING status
// depending on the provided warning duration.
func AdvancedEvaluateDuration(initDuration, serverDuration, duration time.Duration, warning time.Duration) (status Status, message string) {
	if duration >= warning {
		status = StatusWarning
	} else {
		status = StatusOK
	}

	message = fmt.Sprintf("%d ms / %d ms", initDuration.Nanoseconds()/1000000, serverDuration.Nanoseconds()/1000000)

	return
}
