package compute

import "fmt"

func TamperFails(sealed, received string) error {
	if sealed == "" || received == "" {
		return fmt.Errorf("TEE_VERIFY_FAIL")
	}
	if sealed != received {
		return fmt.Errorf("TEE_VERIFY_FAIL")
	}
	return RequireScheme(received)
}
