package exec

import "fmt"

func RejectMutation(bound, presented string) error {
	if bound == "" || presented == "" {
		return fmt.Errorf("preview_incomplete")
	}
	if bound != presented {
		return fmt.Errorf("preview_mutated")
	}
	return nil
}
