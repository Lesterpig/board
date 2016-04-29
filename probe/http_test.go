package probe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPSuccess(t *testing.T) {
	h := NewHTTP("https://www.lesterpig.com", 5*time.Second, 10*time.Second, "Loïck")
	s, m := h.Probe()
	assert.Equal(t, StatusOK, s)
	t.Log(m)
}

func TestHTTPWarning(t *testing.T) {
	h := NewHTTP("https://www.lesterpig.com", time.Millisecond, 10*time.Second, "")
	s, _ := h.Probe()
	assert.Equal(t, StatusWarning, s)
}

func TestHTTP404(t *testing.T) {
	h := NewHTTP("https://www.lesterpig.com/404", 5*time.Second, 10*time.Second, "")
	s, _ := h.Probe()
	assert.Equal(t, StatusError, s)
}

func TestHTTPError(t *testing.T) {
	h := NewHTTP("https://www.leeesterpig.com", 5*time.Second, 10*time.Second, "")
	s, _ := h.Probe()
	assert.Equal(t, StatusError, s)
}

func TestHTTPTimeout(t *testing.T) {
	h := NewHTTP("https://www.lesterpig.com", time.Second, time.Millisecond, "")
	s, _ := h.Probe()
	assert.Equal(t, StatusError, s)
}

func TestHTTPUnexpected(t *testing.T) {
	h := NewHTTP("https://www.lesterpig.com", 5*time.Second, 10*time.Second, "Unexpected")
	s, m := h.Probe()
	assert.Equal(t, StatusError, s)
	assert.Equal(t, "Unexpected result", m)
}
