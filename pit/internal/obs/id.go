package obs

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func NewRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("req_fallback")
	}
	return "req_" + hex.EncodeToString(b[:])
}
