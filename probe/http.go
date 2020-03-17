package probe

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !opts.VerifyCertificate},
	}

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
	start := time.Now()
	res, err := h.client.Get(h.Target)
	duration := time.Since(start)

	if err != nil {
		return StatusError, defaultConnectErrorMsg
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != 200 {
		return StatusError, strconv.Itoa(res.StatusCode)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if h.regex != nil && !h.regex.Match(body) {
		return StatusError, "Unexpected result"
	}

	return EvaluateDuration(duration, h.Warning)
}
