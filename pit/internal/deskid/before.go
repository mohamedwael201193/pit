package deskid

import "fmt"

func BeforeSealedAsk(authorized bool) error {
	if !authorized {
		return fmt.Errorf("desk_not_authorized")
	}
	return nil
}
