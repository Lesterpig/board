package probe

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

// Minecraft Probe, used to check minecraft servers status
type Minecraft struct {
	addrport string
	fatal    time.Duration
}

// NewMinecraft returns a ready-to-go probe.
// The resulting probe will perform a real minecraft handshake to get
// some stats on the server (connected players and version).
func NewMinecraft(addrport string, fatal time.Duration) *Minecraft {
	return &Minecraft{
		addrport: addrport,
		fatal:    fatal,
	}
}

// Probe checks a minecraft server status.
// If the operation succeeds, the message will contain the number of connected
// and allowed players and the server version.
// If there is no slot available for a new player, a warning will be triggered.
// Otherwise, an error message is returned.
func (m *Minecraft) Probe() (status Status, message string) {
	conn, err := net.DialTimeout("tcp", m.addrport, m.fatal)
	if err != nil {
		return StatusError, defaultConnectErrorMsg
	}
	defer func() { _ = conn.Close() }()

	// Handshake
	handshake := []byte{
		0x06, // Length
		0x00, // PacketID
		0x13, // Protocol version varint (109)
		0x00, // String length of server name, not used
		0x00, // Unsigned-short port used, not used
		0x00,
		0x01, // Ask for status
	}
	_, err = conn.Write(handshake)
	if err != nil {
		return StatusError, "Error during handshake"
	}

	// Status
	stat := []byte{0x01, 0x00}
	_, err = conn.Write(stat)
	if err != nil {
		return StatusError, "Error during status"
	}

	// Result
	_, _ = readVarInt(conn) // Packet length
	_, _ = readVarInt(conn) // PacketID
	data := make([]byte, 10000)
	read, err := conn.Read(data)
	if err != nil || read < 2 {
		return StatusError, "No stat received"
	}

	// Try to parse data
	result := new(minecraftServerStats)
	err = json.Unmarshal(data[2:read], result)
	if err != nil {
		return StatusError, "Invalid stats"
	}

	message = fmt.Sprintf("%d / %d - %s", result.Players.Online, result.Players.Max, result.Version.Name)
	status = StatusOK
	if result.Players.Online == result.Players.Max {
		status = StatusWarning
	}
	return
}

type minecraftServerStats struct {
	Version struct {
		Name string `json:"name"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
	} `json:"players"`
}

func readVarInt(c io.Reader) (res int, err error) {
	buf := []byte{0x00}
	for i := 0; i < 5; i++ {
		_, err := c.Read(buf)
		if err != nil {
			return 0, err
		}

		res = res | int((buf[0]&0x7F)<<uint(i*7))
		if 0x00 == buf[0]&0x80 {
			break
		}
	}
	return
}
