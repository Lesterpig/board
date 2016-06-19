package probe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMinecraftSuccess(t *testing.T) {
	h := NewMinecraft("lesterpig.com:25565", 10*time.Second)
	s, m := h.Probe()
	assert.True(t, StatusOK == s)
	t.Log(m)
}

func TestMinecraftError(t *testing.T) {
	h := NewMinecraft("lesterpig.com:25566", 10*time.Second)
	s, _ := h.Probe()
	assert.True(t, StatusError == s)
}
