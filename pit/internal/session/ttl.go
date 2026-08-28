package session

import "fmt"

const (
	DefaultTTLHours = 24
	MaxTTLHours     = 168
)

func CapTTLHours(hours int) error {
	if hours < 1 || hours > MaxTTLHours {
		return fmt.Errorf("session_ttl")
	}
	return nil
}
