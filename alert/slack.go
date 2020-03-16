package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Lesterpig/board/probe"
)

// Slack alert container.
type Slack struct {
	client     *http.Client
	webhookURL string
	channel    string
}

// NewSlack returns a Slack alerter from the webhookURL
func NewSlack(webhookURL, channel string) *Slack {
	return &Slack{
		client:     &http.Client{},
		webhookURL: webhookURL,
		channel:    channel,
	}
}

// Alert sends a pushbullet note to the owner of the provided token.
func (p *Slack) Alert(status probe.Status, category, title, message, link, date string) {

	color := "#ff0000"
	if status == probe.StatusOK {
		color = "#00ff00"
	}
	rocketChat := SlackPostMessage{Channel: p.channel, Alias: category}
	rocketChat.Attachments = []Attachment{
		Attachment{
			Color: color,
			Fields: []AttachmentField{
				{
					Short: false,
					Title: fmt.Sprintf("%s is %s", title, status),
					Value: fmt.Sprintf("[%s](%s) - Response: %s at (%s)", title, link, message, date),
				},
			},
		},
	}

	slackBody, _ := json.Marshal(rocketChat)
	req, err := http.NewRequest(http.MethodPost, p.webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending request to  Slack: %s ", err)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "{\"success\":true}" {
		log.Errorf("Non-ok response returned from Slack: %s ", buf.String())
		return // errors.New("Non-ok response returned from Slack")
	}
	log.Info("Alert sendt")

	return
}

// All models with inspiration from https://github.com/RocketChat/Rocket.Chat.Go.SDK/blob/master/models/message.go

// SlackPostMessage is the main model for sending messages
type SlackPostMessage struct {
	RoomID      string       `json:"roomId,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	ParseUrls   bool         `json:"parseUrls,omitempty"`
	Alias       string       `json:"alias,omitempty"`
	Emoji       string       `json:"emoji,omitempty"`
	Avatar      string       `json:"avatar,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment Payload for postmessage rest API
//
// https://rocket.chat/docs/developer-guides/rest-api/chat/postmessage/
type Attachment struct {
	Color  string            `json:"color,omitempty"`
	Fields []AttachmentField `json:"fields,omitempty"`
}

// AttachmentField Payload for postmessage rest API
//
// https://rocket.chat/docs/developer-guides/rest-api/chat/postmessage/
type AttachmentField struct {
	Short bool   `json:"short"`
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
}
