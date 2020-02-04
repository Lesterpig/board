package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

// Matrix alert container.
type Matrix struct {
	client *http.Client
	token  string
	room   string
}

// NewMatrix returns a Matrix alerter from the private token
// and the room identifier, both available in the Riot web UI.
func NewMatrix(token string, room string) *Matrix {
	return &Matrix{
		client: &http.Client{},
		token:  token,
		room:   room,
	}
}

// Alert sends a message to the matrix room.
func (p *Matrix) Alert(title, body string) {
	data := "☢️ **" + title + "**: " + body
	markdownBody := blackfriday.Run(
		[]byte(data),
		blackfriday.WithExtensions(blackfriday.HardLineBreak|blackfriday.NoEmptyLineBeforeBlock),
	)
	raw, _ := json.Marshal(map[string]string{
		"msgtype":        "m.text",
		"format":         "org.matrix.custom.html",
		"formatted_body": string(markdownBody),
		"body":           data,
	})

	url := fmt.Sprintf(
		"https://matrix.org/_matrix/client/r0/rooms/%s/send/m.room.message/%d?access_token=%s",
		url.PathEscape(p.room),
		time.Now().UnixNano(),
		p.token,
	)

	req, _ := http.NewRequest("PUT", url, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	_, _ = p.client.Do(req)
}
