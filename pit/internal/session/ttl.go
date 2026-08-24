package session

import "fmt"

func CapTTLHours(hours int) error {
	if hours < 1 || hours > 1 {
		return fmt.Errorf("session_ttl")
	}
	return nil
}
