package compute

import "fmt"

func HonestLabel(n Independence) string {
	switch n {
	case Providers:
		return "independent providers"
	case Models:
		return "independent models"
	default:
		return "role separation and envelope separation on the same provider"
	}
}

func RefuseProviderSpoof(claimed, actual string) error {
	if claimed == "" || actual == "" {
		return fmt.Errorf("provider_spoof")
	}
	if claimed != actual {
		return fmt.Errorf("provider_spoof")
	}
	return nil
}
