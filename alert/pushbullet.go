package alert

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Pushbullet alert container.
type Pushbullet struct {
	client *http.Client
	token  string
}

// NewPushbullet returns a Pushbullet alerter from the private token
// available in the `account` page.
func NewPushbullet(token string) *Pushbullet {
	return &Pushbullet{
		client: &http.Client{},
		token:  token,
	}
}

// Alert sends a pushbullet note to the owner of the provided token.
func (p *Pushbullet) Alert(title, body string) {
	u, _ := url.Parse("https://api.pushbullet.com/v2/pushes")
	r := strings.NewReader(`{
	"title": "` + strings.Replace(title, "\"", "\\\"", -1) + `",
	"body": "` + strings.Replace(body, "\"", "\\\"", -1) + `",
	"type": "note"
}`)

	_, _ = p.client.Do(&http.Request{
		Method: "POST",
		URL:    u,
		Header: map[string][]string{
			"Access-Token": {p.token},
			"Content-Type": {"application/json"},
		},
		Body: ioutil.NopCloser(r),
	})
}
