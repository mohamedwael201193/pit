package session

import "fmt"

func CheckTTL(ttlSeconds, maxSeconds int64) error {
	if ttlSeconds <= 0 {
		return fmt.Errorf("ttl_required")
	}
	if maxSeconds > 0 && ttlSeconds > maxSeconds {
		return fmt.Errorf("ttl_too_long")
	}
	return nil
}
