package probe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPSuccess(t *testing.T) {
	h := &HTTP{}
	assert.Nil(t, h.Init(Config{
		Target:  "https://www.lesterpig.com",
		Warning: 5 * time.Second,
		Fatal:   10 * time.Second,
		Options: map[string]interface{}{"Regex": "Loïck"},
	}))

	s, m := h.Probe()

	assert.True(t, StatusOK == s)
	t.Log(m)
}

func TestHTTPWarning(t *testing.T) {
	h := &HTTP{}
	assert.Nil(t, h.Init(Config{
		Target:  "https://www.lesterpig.com",
		Warning: time.Microsecond,
		Fatal:   10 * time.Second,
	}))

	s, _ := h.Probe()

	assert.True(t, StatusWarning == s)
}

func TestHTTP404(t *testing.T) {
	h := &HTTP{}
	assert.Nil(t, h.Init(Config{
		Target:  "https://www.lesterpig.com/404",
		Warning: 5 * time.Second,
		Fatal:   10 * time.Second,
	}))

	s, _ := h.Probe()

	assert.True(t, StatusError == s)
}

func TestHTTPError(t *testing.T) {
	h := &HTTP{}
	assert.Nil(t, h.Init(Config{
		Target:  "https://www.lesteerpig.com",
		Warning: 5 * time.Second,
		Fatal:   10 * time.Second,
	}))

	s, _ := h.Probe()

	assert.True(t, StatusError == s)
}

func TestHTTPTimeout(t *testing.T) {
	h := &HTTP{}
	assert.Nil(t, h.Init(Config{
		Target:  "https://www.lesterpig.com/",
		Warning: 5 * time.Second,
		Fatal:   time.Microsecond,
	}))

	s, _ := h.Probe()

	assert.True(t, StatusError == s)
}

func TestHTTPUnexpected(t *testing.T) {
	h := &HTTP{}
	assert.Nil(t, h.Init(Config{
		Target:  "https://www.lesterpig.com",
		Warning: 5 * time.Second,
		Fatal:   10 * time.Second,
		Options: map[string]interface{}{"Regex": "Unexpected"},
	}))

	s, m := h.Probe()

	assert.True(t, StatusError == s)
	assert.True(t, m == "Unexpected result")
}
