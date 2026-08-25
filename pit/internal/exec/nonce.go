package exec

import "fmt"

func BindNonce(preview, presented int64) error {
	if preview == 0 || presented == 0 || preview != presented {
		return fmt.Errorf("nonce_mismatch")
	}
	return nil
}
