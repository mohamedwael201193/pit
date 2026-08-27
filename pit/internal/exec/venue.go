package exec

import "fmt"

func NeedOnVenue(found bool) error {
	if !found {
		return fmt.Errorf("not_on_venue")
	}
	return nil
}
