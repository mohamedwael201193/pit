package cli

import "fmt"

func RefuseNetworkSwitch(bound, want string) error {
	if bound == "" || want == "" {
		return fmt.Errorf("network_required")
	}
	if bound != want {
		return fmt.Errorf("network_switch_denied")
	}
	return nil
}
