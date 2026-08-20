package identity

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Address is a checksum-insensitive EVM address. Stored lowercase.
type Address string

func NormalizeAddress(raw string) (Address, error) {
	s := strings.TrimSpace(raw)
	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return "", fmt.Errorf("not an address")
	}
	hex := s[2:]
	if len(hex) != 40 {
		return "", fmt.Errorf("address length")
	}
	for _, c := range hex {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return "", fmt.Errorf("address hex")
		}
	}
	return Address(strings.ToLower("0x" + hex)), nil
}

func NewWorkspaceID() string {
	return uuid.NewString()
}

func ParseWorkspaceID(s string) (string, error) {
	id, err := uuid.Parse(strings.TrimSpace(s))
	if err != nil {
		return "", fmt.Errorf("not found")
	}
	if id.Version() != 4 {
		return "", fmt.Errorf("not found")
	}
	return id.String(), nil
}
