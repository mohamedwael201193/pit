package exec

import "fmt"

func DuplicateClick(applied bool, err error) error {
	if err != nil {
		return err
	}
	if !applied {
		return fmt.Errorf("duplicate_click")
	}
	return nil
}
