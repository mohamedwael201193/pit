package deskid

import "fmt"

func TestnetDeskRequired(addr string, mint bool) error {
	if mint && addr == "" {
		return fmt.Errorf("desk_not_deployed")
	}
	return nil
}
